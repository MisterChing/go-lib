package netx

import (
	"net"
)

var localIP = ""

var LanIPSeg = [4]string{
	"10.0.0.0/8",
	"172.16.0.0/12",
	"192.168.0.0/16",
	"127.0.0.1/8",
}

func init() {
	addrs, _ := net.InterfaceAddrs()
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				localIP = ipnet.IP.String()
				return
			}
		}
	}
}

// IsLan 判断是否是局域网ip
func IsLan(ipStr string) bool {
	if ip := net.ParseIP(ipStr); ip != nil {
		for _, network := range LanIPSeg {
			_, subnet, _ := net.ParseCIDR(network)
			if subnet.Contains(ip) {
				return true
			}
		}
	}
	return false
}

// GetLocalIP 获取本地IP
func GetLocalIP() string {
	return localIP
}
