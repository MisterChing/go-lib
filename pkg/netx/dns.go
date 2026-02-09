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

// DNSConfig DNS 配置
type DNSConfig struct {
	Name    string        // DNS 配置名称，用于标识不同的 DNS 配置（如 "Google", "Aliyun"）
	Servers []string      // DNS 服务器列表，如 ["8.8.8.8:53", "114.114.114.114:53"]
	Timeout time.Duration // 超时时间
}

// DefaultDNSConfig 默认配置（使用常见的公共 DNS）
var DefaultDNSConfig = &DNSConfig{
	Name: "Default",
	Servers: []string{
		"8.8.8.8:53",         // Google DNS
		"114.114.114.114:53", // 114 DNS
		"223.5.5.5:53",       // 阿里 DNS
	},
	Timeout: 5 * time.Second,
}

// DNSQueryResult DNS 查询结果
// 支持多种查询方式：DNS 服务器查询、本地 hosts 文件解析等
type DNSQueryResult struct {
	ConfigName string // DNS 配置名称（如 "Default", "Google", "Local Hosts"）
	Server     string // 实际使用的 DNS 服务器（如 "8.8.8.8:53", "file:///etc/hosts"）
	Error      string // 错误信息（如果有）
	IP         string // 解析得到的 IP
}

// SearchIpByDomainWithMultipleDNS 使用多个 DNS 配置查询同一个域名，返回所有结果（用于对比不同 DNS 的解析结果）
// 默认配置始终会被追加到第一个位置
func SearchIpByDomainWithMultipleDNS(domain string, configs []*DNSConfig) []*DNSQueryResult {
	// 始终将默认配置放在第一个位置
	allConfigs := make([]*DNSConfig, 0, len(configs)+1)
	allConfigs = append(allConfigs, DefaultDNSConfig)

	if len(configs) > 0 {
		allConfigs = append(allConfigs, configs...)
	}

	// 预分配结果数组，保持与 allConfigs 相同的顺序
	results := make([]*DNSQueryResult, len(allConfigs))

	// 使用带索引的 channel 来保持顺序
	type indexedResult struct {
		index  int
		result *DNSQueryResult
	}
	resultCh := make(chan indexedResult, len(allConfigs))

	// 并发查询所有 DNS 配置（传递索引以保持顺序）
	for i, config := range allConfigs {
		go func(idx int, cfg *DNSConfig) {
			result := &DNSQueryResult{
				ConfigName: cfg.Name,
			}

			ip, err := searchWithSingleDNSConfig(domain, cfg)
			result.IP = ip
			if err != nil {
				result.Error = err.Error()
			}

			// 记录实际使用的服务器（第一个成功的）
			if err == nil && len(cfg.Servers) > 0 {
				result.Server = cfg.Servers[0]
			}

			// 发送结果时带上索引
			resultCh <- indexedResult{index: idx, result: result}
		}(i, config)
	}

	// 收集所有结果并放到对应的索引位置（保持与 allConfigs 相同的顺序）
	for i := 0; i < len(allConfigs); i++ {
		ir := <-resultCh
		results[ir.index] = ir.result
	}

	return results
}

// searchWithSingleDNSConfig 使用单个 DNS 配置查询（内部辅助函数）
func searchWithSingleDNSConfig(domain string, config *DNSConfig) (string, error) {
	if config == nil {
		config = DefaultDNSConfig
	}

	// 确保域名以 . 结尾
	if domain[len(domain)-1] != '.' {
		domain = domain + "."
	}

	client := &dnsv1.Client{
		Timeout: config.Timeout,
	}

	msg := &dnsv1.Msg{}
	msg.SetQuestion(domain, dnsv1.TypeA)
	msg.RecursionDesired = true

	// 尝试所有配置的 DNS 服务器
	var lastErr error
	for _, server := range config.Servers {
		resp, _, err := client.Exchange(msg, server)
		if err != nil {
			lastErr = fmt.Errorf("query %s failed: %w", server, err)
			continue
		}

		if resp == nil || resp.Rcode != dnsv1.RcodeSuccess {
			lastErr = fmt.Errorf("query %s failed: invalid response", server)
			continue
		}

		// 提取 A 记录
		for _, answer := range resp.Answer {
			if a, ok := answer.(*dnsv1.A); ok {
				return a.A.String(), nil
			}
		}

		lastErr = fmt.Errorf("no A record found in response from %s", server)
	}

	if lastErr != nil {
		return "", lastErr
	}
	return "", errors.New("no ip found")
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
