package sign

import (
	"crypto/hmac"
	"crypto/sha1"
)

func HmacSign(secretKey, body string) []byte {
	m := hmac.New(sha1.New, []byte(secretKey))
	m.Write([]byte(body))
	return m.Sum(nil)
}
