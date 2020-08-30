package utils

import (
	"strings"

	"github.com/qiniu/api.v7/storage"
	"github.com/spf13/viper"

	"github.com/1024casts/snake/pkg/constvar"
)

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

// GetQiNiuPublicAccessURL 获取七牛资源的公有链接
// 无需配置bucket, 域名会自动到域名所绑定的bucket去查找
func GetQiNiuPublicAccessURL(path string) string {
	domain := viper.GetString("qiniu.cdn_url")
	key := strings.TrimPrefix(path, "/")

	publicAccessURL := storage.MakePublicURL(domain, key)

	return publicAccessURL
}
