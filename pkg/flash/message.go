package flash

// https://www.alexedwards.net/blog/simple-flash-messages-in-golang

import (
	"encoding/base64"
	"net/http"
	"time"

	"github.com/1024casts/snake/pkg/log"
)

var flashName = "flash"

func HasFlash(r *http.Request) bool {
	c, err := r.Cookie(flashName)
	if err != nil {
		switch err {
		case http.ErrNoCookie:
			return false
		default:
			return false
		}
	}
	log.Warnf("[flash] read cookie err: %v", err)
	if c.Value != "" {
		return true
	}

	return false
}

func SetMessage(w http.ResponseWriter, msg string) {
	log.Info("[flash] begin set message...")
	expire := time.Now().Add(3 * time.Second)
	value := []byte(msg)
	c := http.Cookie{
		Name:    flashName,
		Value:   base64.URLEncoding.EncodeToString(value),
		Path:    "/",
		Expires: expire,
		MaxAge:  3,
	}
	http.SetCookie(w, &c)
}

func GetMessage(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	log.Info("[flash] begin get message...")
	c, err := r.Cookie(flashName)
	if err != nil {
		switch err {
		case http.ErrNoCookie:
			return nil, nil
		default:
			return nil, err
		}
	}

	// delete cookie
	dc := http.Cookie{
		Name:    flashName,
		Path:    "/",
		MaxAge:  -1,
		Expires: time.Unix(1, 0),
	}
	http.SetCookie(w, &dc)

	value, err := base64.URLEncoding.DecodeString(c.Value)
	if err != nil {
		return nil, err
	}

	return value, nil
}
