package sign

import "crypto/md5"

// Md5Sign md5 sign
func Md5Sign(_, body string) []byte {
	m := md5.New()
	_, _ = m.Write([]byte(body))
	return m.Sum(nil)
}
