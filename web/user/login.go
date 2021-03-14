package user

import (
	"net/http"

	"github.com/1024casts/snake/internal/ecode"

	"github.com/gin-gonic/gin"

	"github.com/1024casts/snake/internal/service"
	"github.com/1024casts/snake/pkg/errno"
	"github.com/1024casts/snake/pkg/log"
	"github.com/1024casts/snake/web"
)

// GetLogin show login page
func GetLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "user/login", gin.H{
		"title": "登录",
		"ctx":   c,
	})
}

// DoLogin perform login
func DoLogin(c *gin.Context) {
	log.Info("[web.login] User DoLogin function called.")
	// Binding the data with the user struct.
	var req LoginCredentials
	if err := c.Bind(&req); err != nil {
		log.Warnf("[web.login] bind err: %v", err)
		web.Response(c, errno.ErrBind, nil)
		return
	}

	// Get the user information by the login username.
	d, err := service.UserSvc.GetUserByEmail(c, req.Email)
	if err != nil {
		log.Warnf("[web.login] get user by email err: %v", err)
		web.Response(c, ecode.ErrUserNotFound, nil)
		return
	}

	log.Info("userbase", d.Password)
	log.Info("req", req.Password)
	// Compare the login password with the user password.
	if err := d.Compare(req.Password); err != nil {
		log.Warnf("[web.login] compare user password err: %v", err)
		web.Response(c, ecode.ErrPasswordIncorrect, nil)
		return
	}

	// set cookie 30 day
	web.SetLoginCookie(c, d.ID)

	web.Response(c, nil, nil)
}
