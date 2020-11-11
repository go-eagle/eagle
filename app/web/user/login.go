package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/1024casts/snake/app/web"
	"github.com/1024casts/snake/internal/service"
	"github.com/1024casts/snake/pkg/errno"
	"github.com/1024casts/snake/pkg/log"
)

// GetLogin show login page
func GetLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "user/login", gin.H{
		"title": "登录",
		"ctx":   c,
	})
}

// DoLogin perform login
// @Summary Login generates the authentication token
// @Produce  json
// @Param username body string true "Username"
// @Param password body string true "Password"
// @Success 200 {string} json "{"code":0,"message":"OK","data":{"token":"eyJhbGciOiJIUzI1NiIsInR7HXsSLfebw92DKd1JQ"}}"
// @Router /login [post]
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
	d, err := service.Svc.UserSvc().GetUserByEmail(c, req.Email)
	if err != nil {
		log.Warnf("[web.login] get user by email err: %v", err)
		web.Response(c, errno.ErrUserNotFound, nil)
		return
	}

	log.Info("userbase", d.Password)
	log.Info("req", req.Password)
	// Compare the login password with the user password.
	if err := d.Compare(req.Password); err != nil {
		log.Warnf("[web.login] compare user password err: %v", err)
		web.Response(c, errno.ErrPasswordIncorrect, nil)
		return
	}

	// set cookie 30 day
	web.SetLoginCookie(c, d.ID)

	web.Response(c, nil, nil)
}
