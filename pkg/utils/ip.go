package utils

import (
	"net"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	tnet "github.com/toolkits/net"
)

var (
	once     sync.Once
	clientIP = "127.0.0.1"
)

// GetLocalIP 获取本地内网IP
func GetLocalIP() string {
	once.Do(func() {
		ips, _ := tnet.IntranetIP()
		if len(ips) > 0 {
			clientIP = ips[0]
		} else {
			clientIP = "127.0.0.1"
		}
	})
	return clientIP
}

// GetInternalIP get internal ip.
func GetInternalIP() string {
	inters, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, inter := range inters {
		if inter.Flags&net.FlagUp != 0 && !strings.HasPrefix(inter.Name, "lo") {
			addrs, err := inter.Addrs()
			if err != nil {
				continue
			}
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return ipnet.IP.String()
					}
				}
			}
		}
	}
	return ""
}

// GetRealIP get user real ip
func GetRealIP(ctx *gin.Context) (ip string) {
	var header = ctx.Request.Header
	var index int
	if ip = header.Get("X-Forwarded-For"); ip != "" {
		index = strings.IndexByte(ip, ',')
		if index < 0 {
			return ip
		}
		if ip = ip[:index]; ip != "" {
			return ip
		}
	}
	if ip = header.Get("X-Real-Ip"); ip != "" {
		index = strings.IndexByte(ip, ',')
		if index < 0 {
			return ip
		}
		if ip = ip[:index]; ip != "" {
			return ip
		}
	}
	if ip = header.Get("Proxy-Forwarded-For"); ip != "" {
		index = strings.IndexByte(ip, ',')
		if index < 0 {
			return ip
		}
		if ip = ip[:index]; ip != "" {
			return ip
		}
	}
	ip, _, _ = net.SplitHostPort(ctx.Request.RemoteAddr)
	return ip
}
