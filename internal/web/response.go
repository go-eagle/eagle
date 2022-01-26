package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/spf13/viper"

	"github.com/go-eagle/eagle/pkg/errcode"
	"github.com/go-eagle/eagle/pkg/flash"
	"github.com/go-eagle/eagle/pkg/log"
)

// Resp web response struct
type Resp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Response json response
func Response(c *gin.Context, err error, data interface{}) {
	code, message := errcode.DecodeErr(err)

	// always return http.StatusOK
	c.JSON(http.StatusOK, Resp{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// Redirect redirect to another path
func Redirect(c *gin.Context, redirectPath, errMsg string) {
	flash.SetMessage(c.Writer, errMsg)
	c.Redirect(http.StatusMovedPermanently, redirectPath)
	c.Abort()
}

// Store It is recommended to use an authentication key with 32 or 64 bytes.
// The encryption key, if set, must be either 16, 24, or 32 bytes to select
// AES-128, AES-192, or AES-256 modes.
var Store = sessions.NewCookieStore([]byte(viper.GetString("cookie.secret")))

// SetLoginCookie set login cookie to session
func SetLoginCookie(ctx *gin.Context, userID uint64) {
	Store.Options.HttpOnly = true

	session := GetCookieSession(ctx)
	session.Options = &sessions.Options{
		Domain:   "",
		MaxAge:   86400,
		Path:     "/",
		HttpOnly: true,
	}
	// 浏览器关闭，cookie删除，否则保存30天(github.com/gorilla/sessions 包的默认值)
	val := ctx.DefaultPostForm("remember_me", "0")
	log.Infof("[handler] remember_me val: %s", val)
	if val == "1" {
		session.Options.MaxAge = 86400 * 30
	}

	session.Values["user_id"] = userID

	req := Request(ctx)
	resp := ResponseWriter(ctx)
	err := session.Save(req, resp)
	if err != nil {
		log.Warnf("[handler] set login cookie, %v", err)
	}
}

// GetCookieSession get cookie
func GetCookieSession(ctx *gin.Context) *sessions.Session {
	session, err := Store.Get(ctx.Request, "")
	if err != nil {
		log.Warnf("[handler] store get err, %v", err)
	}
	return session
}

// Request return a request
func Request(ctx *gin.Context) *http.Request {
	return ctx.Request
}

// ResponseWriter return a response writer
func ResponseWriter(ctx *gin.Context) http.ResponseWriter {
	return ctx.Writer
}
