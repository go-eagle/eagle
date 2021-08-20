package sign

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/go-eagle/eagle/pkg/utils"
)

// Verifier define struct
type Verifier struct {
	*DefaultKeyName
	body url.Values

	timeout time.Duration // 签名过期时间
}

// NewVerifier 实例化 Verifier
func NewVerifier() *Verifier {
	return &Verifier{
		DefaultKeyName: newDefaultKeyName(),
		body:           make(url.Values),
		timeout:        time.Minute * 5,
	}
}

// ParseQuery 将参数字符串解析成参数列表
func (v *Verifier) ParseQuery(requestURI string) error {
	requestQuery := ""
	idx := strings.Index(requestURI, "?")
	if idx > 0 {
		requestQuery = requestURI[idx+1:]
	}
	query, err := url.ParseQuery(requestQuery)
	if nil != err {
		return err
	}
	v.ParseValues(query)
	return nil
}

// ParseValues 将Values参数列表解析成参数Map。如果参数是多值的，则将它们以逗号Join成字符串。
func (v *Verifier) ParseValues(values url.Values) {
	for key, value := range values {
		v.body[key] = value
	}
}

// SetTimeout 设置签名校验过期时间
func (v *Verifier) SetTimeout(timeout time.Duration) *Verifier {
	v.timeout = timeout
	return v
}

// MustString 获取字符串值
func (v *Verifier) MustString(key string) string {
	ss := v.MustStrings(key)
	if len(ss) == 0 {
		return ""
	}
	return ss[0]
}

// MustStrings 获取字符串值数组
func (v *Verifier) MustStrings(key string) []string {
	return v.body[key]
}

// MustInt64 获取Int64值
func (v *Verifier) MustInt64(key string) int64 {
	n, _ := utils.StringToInt64(v.MustString(key))
	return n
}

// MustHasKeys 必须包含指定的字段参数
func (v *Verifier) MustHasKeys(keys ...string) error {
	for _, key := range keys {
		if _, hit := v.body[key]; !hit {
			return fmt.Errorf("KEY_MISSED:<%s>", key)
		}
	}
	return nil
}

// MustHasOtherKeys 必须包含除特定的[timestamp, nonce_str, sign, app_id]等之外的指定的字段参数
func (v *Verifier) MustHasOtherKeys(keys ...string) error {
	fields := []string{v.Timestamp, v.NonceStr, v.Sign, v.AppID}
	if len(keys) > 0 {
		fields = append(fields, keys...)
	}
	return v.MustHasKeys(fields...)
}

// CheckTimeStamp 检查时间戳有效期
func (v *Verifier) CheckTimeStamp() error {
	timestamp := v.GetTimestamp()
	thatTime := time.Unix(timestamp, 0)
	if timestamp > time.Now().Unix() || time.Since(thatTime) > v.timeout {
		return fmt.Errorf("TIMESTAMP_TIMEOUT:<%d>", timestamp)
	}
	return nil
}

// GetAppID 获取app id
func (v *Verifier) GetAppID() string {
	return v.MustString(v.AppID)
}

// GetNonceStr 获取随机字符串
func (v *Verifier) GetNonceStr() string {
	return v.MustString(v.NonceStr)
}

// GetSign 获取签名
func (v *Verifier) GetSign() string {
	return v.MustString(v.Sign)
}

// GetTimestamp 获取时间戳
func (v *Verifier) GetTimestamp() int64 {
	return v.MustInt64(v.Timestamp)
}

// GetBodyWithoutSign 获取所有参数体。其中不包含sign 字段
func (v *Verifier) GetBodyWithoutSign() url.Values {
	out := make(url.Values)
	for k, val := range v.body {
		if k != v.Sign {
			out[k] = val
		}
	}
	return out
}

// GetBody 获取body
func (v *Verifier) GetBody() url.Values {
	out := make(url.Values)
	for k, val := range v.body {
		out[k] = val
	}
	return out
}
