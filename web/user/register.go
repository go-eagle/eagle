package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/1024casts/snake/internal/model"
	"github.com/1024casts/snake/internal/service"
	"github.com/1024casts/snake/pkg/errno"
	"github.com/1024casts/snake/pkg/flash"
	"github.com/1024casts/snake/pkg/log"
	"github.com/1024casts/snake/web"
)

// GetRegister register as a new user
func GetRegister(c *gin.Context) {
	c.HTML(http.StatusOK, "user/register", gin.H{
		"title": "注册",
		"ctx":   c,
	})
}

// DoRegister submit register
func DoRegister(c *gin.Context) {
	log.Info("User Register function called.")
	var r RegisterRequest
	if err := c.Bind(&r); err != nil {
		web.Response(c, errno.ErrBind, nil)
		return
	}

	u := model.UserBaseModel{
		Username: r.Username,
		Email:    r.Email,
		Password: r.Password,
	}

	// Validate the data.
	if err := u.Validate(); err != nil {
		web.Response(c, errno.ErrValidation, nil)
		return
	}

	// Encrypt the user password.
	if err := u.Encrypt(); err != nil {
		web.Response(c, errno.ErrEncrypt, nil)
		return
	}
	// Insert the user to the database.
	err := service.UserSvc.Register(c, u.Username, u.Email, u.Password)
	if err != nil {
		web.Response(c, errno.InternalServerError, nil)
		return
	}

	flash.SetMessage(c.Writer, "已发送激活链接,请检查您的邮箱。")

	// Show the user information.
	web.Response(c, nil, RegisterResponse{
		ID: u.ID,
	})
}
