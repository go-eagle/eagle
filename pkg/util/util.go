package util

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/teris-io/shortid"
	tnet "github.com/toolkits/net"
)

func GenShortId() (string, error) {
	return shortid.Generate()
}

func GenUUID() string {
	u, _ := uuid.NewRandom()
	return u.String()
}

func GetReqID(c *gin.Context) string {
	v, ok := c.Get("X-Request-Id")
	if !ok {
		return ""
	}
	if requestId, ok := v.(string); ok {
		return requestId
	}
	return ""
}

var (
	once     sync.Once
	clientIP = "127.0.0.1"
)

//GetLocalIP 获取本地内网IP
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
