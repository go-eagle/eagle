package utils

import (
	"strings"

	"github.com/qiniu/api.v7/storage"
)

// GetDefaultAvatarURL 获取默认头像
func GetDefaultAvatarURL(cdnURL string) string {
	uri := "/default/avatar.jpg"
	return GetQiNiuPublicAccessURL(cdnURL, uri)
}

// GetAvatarURL user's avatar, if empty, use default avatar
func GetAvatarURL(cdnURL, key string) string {
	if key == "" {
		return GetDefaultAvatarURL(cdnURL)
	}
	if strings.HasPrefix(key, "https://") {
		return key
	}
	return GetQiNiuPublicAccessURL(cdnURL, key)
}

// GetQiNiuPublicAccessURL 获取七牛资源的公有链接
// 无需配置bucket, 域名会自动到域名所绑定的bucket去查找
func GetQiNiuPublicAccessURL(cdnURL, path string) string {
	domain := cdnURL
	key := strings.TrimPrefix(path, "/")

	publicAccessURL := storage.MakePublicURL(domain, key)

	return publicAccessURL
}
