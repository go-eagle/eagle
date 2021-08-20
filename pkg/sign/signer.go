package sign

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/go-eagle/eagle/pkg/utils"
)

// CryptoFunc 签名加密函数
type CryptoFunc func(secretKey string, args string) []byte

// Signer define
type Signer struct {
	*DefaultKeyName

	body       url.Values // 签名参数体
	bodyPrefix string     // 参数体前缀
	bodySuffix string     // 参数体后缀
	splitChar  string     // 前缀、后缀分隔符号

	secretKey  string // 签名密钥
	cryptoFunc CryptoFunc
}

// NewSigner 实例化 Signer
func NewSigner(cryptoFunc CryptoFunc) *Signer {
	return &Signer{
		DefaultKeyName: newDefaultKeyName(),
		body:           make(url.Values),
		bodyPrefix:     "",
		bodySuffix:     "",
		splitChar:      "",
		cryptoFunc:     cryptoFunc,
	}
}

// NewSignerMd5 md5加密算法
func NewSignerMd5() *Signer {
	return NewSigner(Md5Sign)
}

// NewSignerHmac hmac加密算法
func NewSignerHmac() *Signer {
	return NewSigner(HmacSign)
}

// NewSignerAes aes对称加密算法
func NewSignerAes() *Signer {
	return NewSigner(AesSign)
}

// SetBody 设置整个参数体Body对象。
func (s *Signer) SetBody(body url.Values) {
	for k, v := range body {
		s.body[k] = v
	}
}

// GetBody 返回Body内容
func (s *Signer) GetBody() url.Values {
	return s.body
}

// AddBody 添加签名体字段和值
func (s *Signer) AddBody(key string, value string) *Signer {
	return s.AddBodies(key, []string{value})
}

// AddBodies add value to body
func (s *Signer) AddBodies(key string, value []string) *Signer {
	s.body[key] = value
	return s
}

// SetTimeStamp 设置时间戳参数
func (s *Signer) SetTimeStamp(ts int64) *Signer {
	return s.AddBody(s.Timestamp, strconv.FormatInt(ts, 10))
}

// GetTimeStamp 获取TimeStamp
func (s *Signer) GetTimeStamp() string {
	return s.body.Get(s.Timestamp)
}

// SetNonceStr 设置随机字符串参数
func (s *Signer) SetNonceStr(nonce string) *Signer {
	return s.AddBody(s.NonceStr, nonce)
}

// GetNonceStr 返回NonceStr字符串
func (s *Signer) GetNonceStr() string {
	return s.body.Get(s.NonceStr)
}

// SetAppID 设置AppId参数
func (s *Signer) SetAppID(appID string) *Signer {
	return s.AddBody(s.AppID, appID)
}

// GetAppID get app id
func (s *Signer) GetAppID() string {
	return s.body.Get(s.AppID)
}

// RandNonceStr 自动生成16位随机字符串参数
func (s *Signer) RandNonceStr() *Signer {
	return s.SetNonceStr(utils.RandomStr(16))
}

// SetSignBodyPrefix 设置签名字符串的前缀字符串
func (s *Signer) SetSignBodyPrefix(prefix string) *Signer {
	s.bodyPrefix = prefix
	return s
}

// SetSignBodySuffix 设置签名字符串的后缀字符串
func (s *Signer) SetSignBodySuffix(suffix string) *Signer {
	s.bodySuffix = suffix
	return s
}

// SetSplitChar 设置前缀、后缀与签名体之间的分隔符号。默认为空字符串
func (s *Signer) SetSplitChar(split string) *Signer {
	s.splitChar = split
	return s
}

// SetAppSecret 设置签名密钥
func (s *Signer) SetAppSecret(appSecret string) *Signer {
	s.secretKey = appSecret
	return s
}

// SetAppSecretWrapBody 在签名参数体的首部和尾部，拼接AppSecret字符串。
func (s *Signer) SetAppSecretWrapBody(appSecret string) *Signer {
	s.SetSignBodyPrefix(appSecret)
	s.SetSignBodySuffix(appSecret)
	return s.SetAppSecret(appSecret)
}

// GetSignBodyString 获取用于签名的原始字符串
func (s *Signer) GetSignBodyString() string {
	return s.MakeRawBodyString()
}

// MakeRawBodyString 获取用于签名的原始字符串
func (s *Signer) MakeRawBodyString() string {
	return s.bodyPrefix + s.splitChar + s.getSortedBodyString() + s.splitChar + s.bodySuffix
}

// GetSignedQuery 获取带签名参数的查询字符串
func (s *Signer) GetSignedQuery() string {
	return s.MakeSignedQuery()
}

// MakeSignedQuery 获取带签名参数的字符串
func (s *Signer) MakeSignedQuery() string {
	body := s.getSortedBodyString()
	sign := s.GetSignature()
	return body + "&" + s.Sign + "=" + sign
}

// GetSignature 获取签名
func (s *Signer) GetSignature() string {
	return s.MakeSign()
}

// MakeSign 生成签名
func (s *Signer) MakeSign() string {
	sign := fmt.Sprintf("%x", s.cryptoFunc(s.secretKey, s.GetSignBodyString()))
	return sign
}

func (s *Signer) getSortedBodyString() string {
	return SortKVPairs(s.body)
}

// SortKVPairs 将Map的键值对，按字典顺序拼接成字符串
func SortKVPairs(m url.Values) string {
	size := len(m)
	if size == 0 {
		return ""
	}
	keys := make([]string, size)
	idx := 0
	for k := range m {
		keys[idx] = k
		idx++
	}
	sort.Strings(keys)
	pairs := make([]string, size)
	for i, key := range keys {
		pairs[i] = key + "=" + strings.Join(m[key], ",")
	}
	return strings.Join(pairs, "&")
}
