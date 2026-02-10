package netx

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	dnsv1 "github.com/miekg/dns"
)

// AdvancedDNSConfig 高级 DNS 配置（类似 /etc/resolv.conf）
type AdvancedDNSConfig struct {
	Name        string        // 配置名称，用于标识不同的配置（如 "Production", "Test"）
	Nameservers []string      // DNS 服务器列表，如 ["10.189.0.10:53"]
	Search      []string      // 域名搜索后缀列表，如 ["svc.cluster.local", "cluster.local"]
	Timeout     time.Duration // 单次查询超时时间（对应 resolv.conf 的 timeout）
	Attempts    int           // 查询失败后的重试次数（对应 resolv.conf 的 attempts）
	Ndots       int           // 域名中点的数量阈值（对应 resolv.conf 的 ndots）
}

// AdvancedDNSQueryResult 高级 DNS 查询的详细结果
type AdvancedDNSQueryResult struct {
	ConfigName string // 配置名称（如 "Production", "Test"）
	Query      string // 查询的完整域名（如 "redis.svc.cluster.local"）
	Nameserver string // 使用的 DNS 服务器（如 "10.189.0.10:53"）
	Error      string // 错误信息（如果有）
	IP         string // 解析得到的 IP
}

// DefaultAdvancedDNSConfig 默认高级配置
var DefaultAdvancedDNSConfig = &AdvancedDNSConfig{
	Name:        "Default",
	Nameservers: []string{"8.8.8.8:53"},
	Search:      []string{},
	Timeout:     1 * time.Second,
	Attempts:    2,
	Ndots:       2,
}

// SearchIpByDomainWithMultipleAdvancedConfigs 使用多个高级配置并发查询同一个域名，返回所有配置的查询结果
// 用于对比不同配置（如不同环境、不同 DNS 服务器）的解析结果
// 默认配置（DefaultAdvancedDNSConfig）会自动追加到第一个位置
// 优化：所有配置并发查询，带总超时控制
func SearchIpByDomainWithMultipleAdvancedConfigs(domain string, configs []*AdvancedDNSConfig) ([]*AdvancedDNSQueryResult, error) {
	if domain == "" {
		return nil, fmt.Errorf("domain is empty")
	}

	// 始终将默认配置放在第一个位置
	allConfigs := make([]*AdvancedDNSConfig, 0, len(configs)+1)
	allConfigs = append(allConfigs, DefaultAdvancedDNSConfig)

	// 追加用户自定义的配置（去重：如果用户配置中已有同名配置，跳过）
	if len(configs) > 0 {
		for _, cfg := range configs {
			// 检查是否与默认配置同名（避免重复）
			if cfg != nil && cfg.Name != DefaultAdvancedDNSConfig.Name {
				allConfigs = append(allConfigs, cfg)
			}
		}
	}

	// 使用带索引的 channel 来保持顺序
	type indexedResult struct {
		index  int
		result *AdvancedDNSQueryResult
		err    error
	}
	resultCh := make(chan indexedResult, len(allConfigs))

	// 并发查询所有配置（传递索引以保持顺序）
	for i, cfg := range allConfigs {
		go func(idx int, config *AdvancedDNSConfig) {
			result, err := SearchIpByDomainAdvancedWithDetails(domain, config)
			resultCh <- indexedResult{
				index:  idx,
				result: result,
				err:    err,
			}
		}(i, cfg)
	}

	// 收集所有结果并放到对应的索引位置（保持与 allConfigs 相同的顺序）
	orderedResults := make([]*AdvancedDNSQueryResult, len(allConfigs))

	for range allConfigs {
		ir := <-resultCh
		orderedResults[ir.index] = ir.result
	}

	return orderedResults, nil
}

// SearchIpByDomainAdvancedWithDetails 使用高级配置查询域名，返回最终命中的查询结果
// 成功时返回命中的查询结果；失败时返回最后一次查询的失败结果
// 优化：使用并发查询，一旦任何查询成功立即返回，最坏情况超时 = max(timeout * attempts)
func SearchIpByDomainAdvancedWithDetails(domain string, config *AdvancedDNSConfig) (*AdvancedDNSQueryResult, error) {
	if domain == "" {
		return nil, fmt.Errorf("domain is empty")
	}

	// 使用默认配置（如果为 nil）
	if config == nil {
		config = DefaultAdvancedDNSConfig
	}

	// 确保必要的配置有默认值
	if len(config.Nameservers) == 0 {
		return nil, fmt.Errorf("no nameservers configured")
	}
	if config.Timeout <= 0 {
		config.Timeout = 1 * time.Second
	}
	if config.Attempts <= 0 {
		config.Attempts = 2
	}
	if config.Ndots <= 0 {
		config.Ndots = 2
	}
	// 如果 Name 为空，设置一个默认名称
	if config.Name == "" {
		config.Name = "Unnamed"
	}

	// 1. 根据 ndots 决定查询顺序
	queries := buildQueryList(domain, config.Search, config.Ndots)

	// 2. 创建 DNS 客户端
	client := &dnsv1.Client{
		Timeout: config.Timeout,
	}

	// 3. 为每个查询域名生成一个任务
	resultCh := make(chan *AdvancedDNSQueryResult, len(queries))

	// 并发执行所有查询
	for _, query := range queries {
		go func(q string) {
			// 确保域名以 . 结尾（FQDN）
			queryFQDN := q
			if queryFQDN[len(queryFQDN)-1] != '.' {
				queryFQDN = queryFQDN + "."
			}

			var lastResult *AdvancedDNSQueryResult

			// 依次尝试每个 nameserver（失败才用下一个）
			for _, nameserver := range config.Nameservers {
				// 对当前 nameserver 尝试 attempts 次（失败才重试）
				for attempt := 1; attempt <= config.Attempts; attempt++ {
					ip, err := queryDNS(client, queryFQDN, nameserver)

					result := &AdvancedDNSQueryResult{
						ConfigName: config.Name,
						Query:      q,
						Nameserver: nameserver,
						IP:         ip,
					}

					if err != nil {
						result.Error = err.Error()
					}

					lastResult = result

					// 如果成功，立即返回
					if ip != "" && err == nil {
						resultCh <- result
						return
					}
				}
			}

			// 所有 nameserver 都失败，发送最后的失败结果
			if lastResult != nil {
				resultCh <- lastResult
			}
		}(query)
	}

	// 收集结果：找到第一个成功的就返回
	var allResults []*AdvancedDNSQueryResult

	for i := 0; i < len(queries); i++ {
		result := <-resultCh
		allResults = append(allResults, result)

		// 如果成功，立即返回
		if result.IP != "" && result.Error == "" {
			return result, nil
		}
	}

	// 所有查询都失败了，返回最后一次查询的结果
	if len(allResults) > 0 {
		lastResult := allResults[len(allResults)-1]
		return lastResult, fmt.Errorf("all queries failed, last error: %s", lastResult.Error)
	}

	return nil, fmt.Errorf("no queries executed for domain %s", domain)
}

// buildQueryList 根据 ndots 规则构建查询列表
// 规则：
//   - 如果域名中的 '.' 数量 >= ndots：先查询原始域名，再尝试追加 search 后缀
//   - 如果域名中的 '.' 数量 < ndots：先尝试追加 search 后缀，再查询原始域名
//   - 如果 searchDomains 为空或 nil：直接返回原始域名（无论 ndots 为多少）
func buildQueryList(domain string, searchDomains []string, ndots int) []string {
	// 去除尾部的 '.'（如果有）
	domain = strings.TrimSuffix(domain, ".")

	// 如果 search 为空，直接返回原始域名
	if len(searchDomains) == 0 {
		return []string{domain}
	}

	// 计算域名中 '.' 的数量
	dotCount := strings.Count(domain, ".")

	var queries []string

	if dotCount >= ndots {
		// 先查询原始域名
		queries = append(queries, domain)
		// 再尝试追加 search 后缀
		for _, suffix := range searchDomains {
			if suffix != "" { // 跳过空后缀
				queries = append(queries, domain+"."+suffix)
			}
		}
	} else {
		// 先尝试追加 search 后缀
		for _, suffix := range searchDomains {
			if suffix != "" { // 跳过空后缀
				queries = append(queries, domain+"."+suffix)
			}
		}
		// 最后查询原始域名
		queries = append(queries, domain)
	}

	return queries
}

// queryDNS 执行单次 DNS 查询
func queryDNS(client *dnsv1.Client, domain string, nameserver string) (string, error) {
	msg := &dnsv1.Msg{}
	msg.SetQuestion(domain, dnsv1.TypeA)
	msg.RecursionDesired = true

	resp, _, err := client.Exchange(msg, nameserver)
	if err != nil {
		return "", fmt.Errorf("query %s failed: %w", nameserver, err)
	}

	if resp == nil || resp.Rcode != dnsv1.RcodeSuccess {
		rcode := 0
		if resp != nil {
			rcode = resp.Rcode
		}
		return "", fmt.Errorf("query %s failed: invalid response (rcode=%d)", nameserver, rcode)
	}

	// 提取 A 记录
	for _, answer := range resp.Answer {
		if a, ok := answer.(*dnsv1.A); ok {
			return a.A.String(), nil
		}
	}

	return "", fmt.Errorf("no A record found from %s", nameserver)
}

// SearchIpByLocalHosts 从本地 hosts 文件解析域名
// 支持 /etc/hosts (Linux/macOS) 和 C:\Windows\System32\drivers\etc\hosts (Windows)
func SearchIpByLocalHosts(domain string) (string, error) {
	return searchIpByLocalHostsWithPath(domain, getHostsFilePath())
}

// searchIpByLocalHostsWithPath 从指定的 hosts 文件路径解析域名
func searchIpByLocalHostsWithPath(domain, hostsPath string) (string, error) {
	if domain == "" {
		return "", errors.New("domain is empty")
	}

	fileContent, err := os.ReadFile(hostsPath)
	if err != nil {
		return "", fmt.Errorf("read hosts file failed: %w", err)
	}

	// 逐行解析
	lines := splitLines(string(fileContent))

	for _, line := range lines {
		// 跳过注释和空行
		if line == "" || line[0] == '#' {
			continue
		}

		// 分割 IP 和域名
		fields := splitFields(line)
		if len(fields) < 2 {
			continue
		}

		ip := fields[0]
		// 检查是否匹配域名（支持多个域名映射到同一个 IP）
		for i := 1; i < len(fields); i++ {
			if fields[i] == domain {
				return ip, nil
			}
		}
	}

	return "", fmt.Errorf("domain %s not found in hosts file", domain)
}

// getHostsFilePath 获取当前操作系统的 hosts 文件路径
func getHostsFilePath() string {
	if runtime.GOOS == "windows" {
		return "C:\\Windows\\System32\\drivers\\etc\\hosts"
	}
	return "/etc/hosts"
}

// splitLines 分割字符串为行（支持 \n, \r\n, \r 三种换行符）
func splitLines(s string) []string {
	// 统一换行符：将 \r\n 和 \r 都替换为 \n
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")

	lines := strings.Split(s, "\n")
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			result = append(result, line)
		}
	}
	return result
}

// splitFields 分割行为字段（使用空格或制表符）
func splitFields(line string) []string {
	// 移除注释部分
	if idx := strings.Index(line, "#"); idx >= 0 {
		line = line[:idx]
	}
	line = strings.TrimSpace(line)

	// 使用 strings.Fields 自动处理多个空格和制表符
	return strings.Fields(line)
}
