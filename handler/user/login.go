package user

import (
	"github.com/gin-gonic/gin"

	"github.com/1024casts/snake/handler"
	"github.com/1024casts/snake/model"
	"github.com/1024casts/snake/pkg/errno"
	"github.com/1024casts/snake/pkg/log"
	"github.com/1024casts/snake/service/user"
	"github.com/1024casts/snake/service/vcode"
)

// Login 邮箱登录
// @Summary 用户登录接口
// @Description 仅限邮箱登录
// @Tags 用户
// @Produce  json
// @Param req body LoginCredentials true ""
// @Success 200 {string} json "{"code":0,"message":"OK","data":{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6Ik"}}"
// @Router /login [post]
func Login(c *gin.Context) {
	// Binding the data with the u struct.
	var req LoginCredentials
	if err := c.Bind(&req); err != nil {
		log.Warnf("email login bind param err: %v", err)
		handler.SendResponse(c, errno.ErrBind, nil)
		return
	}

	log.Infof("login req %#v", req)
	// check param
	if req.Email == "" || req.Password == "" {
		log.Warnf("email or password is empty: %v", req)
		handler.SendResponse(c, errno.ErrParam, nil)
		return
	}

	t, err := user.Svc.EmailLogin(c, req.Email, req.Password)
	if err != nil {
		log.Warnf("email login err: %v", err)
		handler.SendResponse(c, errno.ErrEmailOrPassword, nil)
		return
	}

	handler.SendResponse(c, nil, model.Token{
		Token: t,
	})
}

// PhoneLogin 手机登录接口
// @Summary 用户登录接口
// @Description 仅限手机登录
// @Tags 用户
// @Produce  json
// @Param req body PhoneLoginCredentials true "phone"
// @Success 200 {string} json "{"code":0,"message":"OK","data":{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6Ik"}}"
// @Router /users/login [post]
func PhoneLogin(c *gin.Context) {
	log.Info("Phone Login function called.")

	// Binding the data with the u struct.
	var req PhoneLoginCredentials
	if err := c.Bind(&req); err != nil {
		log.Warnf("phone login bind param err: %v", err)
		handler.SendResponse(c, errno.ErrBind, nil)
		return
	}

	log.Infof("req %#v", req)
	// check param
	if req.Phone == 0 || req.VerifyCode == 0 {
		log.Warn("phone login bind param is empty")
		handler.SendResponse(c, errno.ErrParam, nil)
		return
	}

	// 验证校验码
	if !vcode.VCodeService.CheckLoginVCode(req.Phone, req.VerifyCode) {
		handler.SendResponse(c, errno.ErrVerifyCode, nil)
		return
	}

	// 登录
	t, err := user.Svc.PhoneLogin(c, req.Phone, req.VerifyCode)
	if err != nil {
		handler.SendResponse(c, errno.ErrVerifyCode, nil)
		return
	}

	handler.SendResponse(c, nil, model.Token{
		Token: t,
	})
}
