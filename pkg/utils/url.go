package utils

import (
	"strings"

	"github.com/qiniu/api.v7/storage"

	"github.com/1024casts/snake/pkg/conf"
)

// GetDefaultAvatarURL 获取默认头像
func GetDefaultAvatarURL() string {
	url := "/default/avatar.jpg"
	return GetQiNiuPublicAccessURL(url)
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
	domain := conf.Conf.QiNiu.CdnURL
	key := strings.TrimPrefix(path, "/")

	publicAccessURL := storage.MakePublicURL(domain, key)

	return publicAccessURL
}
