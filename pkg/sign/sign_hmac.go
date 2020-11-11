package sign

import (
	"crypto/hmac"
	"crypto/sha1"
)

// HmacSign hmac
func HmacSign(secretKey, body string) []byte {
	m := hmac.New(sha1.New, []byte(secretKey))
	_, _ = m.Write([]byte(body))
	return m.Sum(nil)
}
