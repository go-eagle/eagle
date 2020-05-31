package util

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/qiniu/api.v7/storage"
	"github.com/spf13/viper"
	"github.com/teris-io/shortid"
	tnet "github.com/toolkits/net"

	"github.com/1024casts/snake/pkg/constvar"
)

// GenShortID 生成一个id
func GenShortID() (string, error) {
	return shortid.Generate()
}

// GenUUID 生成随机字符串，eg: 76d27e8c-a80e-48c8-ad20-e5562e0f67e4
func GenUUID() string {
	u, _ := uuid.NewRandom()
	return u.String()
}

// GetReqID 获取请求中的request_id
func GetReqID(c *gin.Context) string {
	v, ok := c.Get("X-Request-ID")
	if !ok {
		return ""
	}
	if requestID, ok := v.(string); ok {
		return requestID
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

// GetBytes interface 转 byte
func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GetDate 获取字符串日期
func GetDate() string {
	return time.Now().Format("2006/01/02")
}

// GetTodayDateInt 获取整形的日期
func GetTodayDateInt() int {
	dateStr := time.Now().Format("200601")
	date, err := strconv.Atoi(dateStr)
	if err != nil {
		return 0
	}
	return date
}

// TimeLayout 常用日期格式化模板
func TimeLayout() string {
	return "2006-01-02 15:04:05"
}

// TimeToString 时间转字符串
func TimeToString(ts time.Time) string {
	return time.Unix(ts.Unix(), 00).Format(TimeLayout())
}

// TimeToShortString 时间转日期
func TimeToShortString(ts time.Time) string {
	return time.Unix(ts.Unix(), 00).Format("2006.01.02")
}

// GetDefaultAvatarURL 获取默认头像
func GetDefaultAvatarURL() string {
	return GetQiNiuPublicAccessURL(constvar.DefaultAvatar)
}

// GetAvatarURL user's avatar, if empty, use default avatar
func GetAvatarURL(key string) string {
	if key == "" {
		return GetDefaultAvatarURL()
	}
	if strings.HasPrefix(key, "https://") {
		return key
	}
	return GetQiNiuPublicAccessURL(key)
}

// GetShowTime 格式化时间
func GetShowTime(ts time.Time) string {
	duration := time.Now().Unix() - ts.Unix()
	timeStr := ""
	if duration < 60 {
		timeStr = "刚刚发布"
	} else if duration < 3600 {
		timeStr = fmt.Sprintf("%d分钟前更新", duration/60)
	} else if duration < 86400 {
		timeStr = fmt.Sprintf("%d小时前更新", duration/3600)
	} else if duration < 86400*2 {
		timeStr = "昨天更新"
	} else {
		timeStr = TimeToShortString(ts) + "前更新"
	}
	return timeStr
}

// Md5 字符串转md5
func Md5(str string) (string, error) {
	h := md5.New()

	_, err := io.WriteString(h, str)
	if err != nil {
		return "", err
	}

	// 注意：这里不能使用string将[]byte转为字符串，否则会显示乱码
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// GetQiNiuPublicAccessURL 获取七牛资源的公有链接
// 无需配置bucket, 域名会自动到域名所绑定的bucket去查找
func GetQiNiuPublicAccessURL(path string) string {
	domain := viper.GetString("qiniu.cdn_url")
	key := strings.TrimPrefix(path, "/")

	publicAccessURL := storage.MakePublicURL(domain, key)

	return publicAccessURL
}
