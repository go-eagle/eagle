package user

import (
	"github.com/gin-gonic/gin"

	"github.com/go-eagle/eagle/internal/ecode"
	"github.com/go-eagle/eagle/internal/model"
	"github.com/go-eagle/eagle/internal/service"
	"github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/errcode"
	"github.com/go-eagle/eagle/pkg/log"
)

// Login 邮箱登录
// @Summary 用户登录接口
// @Description 仅限邮箱登录
// @Tags 用户
// @Produce  json
// @Param req body LoginCredentials true "请求参数"
// @Success 200 {object} model.UserInfo "用户信息"
// @Router /login [post]
func Login(c *gin.Context) {
	// Binding the data with the u struct.
	var req LoginCredentials
	valid, errs := app.BindAndValid(c, &req)
	if !valid {
		log.Warnf("app.BindAndValid errs: %v", errs)
		response.Error(c, errcode.ErrInvalidParam.WithDetails(errs.Errors()...))
		return
	}

	log.Infof("login req %#v", req)

	t, err := service.Svc.Users().EmailLogin(c, req.Email, req.Password)
	if err != nil {
		log.Warnf("email login err: %v", err)
		response.Error(c, ecode.ErrEmailOrPassword)
		return
	}

	response.Success(c, model.Token{Token: t})
}

// PhoneLogin 手机登录接口
// @Summary 用户登录接口
// @Description 仅限手机登录
// @Tags 用户
// @Produce  json
// @Param req body PhoneLoginCredentials true "phone"
// @Success 200 {object} model.UserInfo "用户信息"
// @Router /users/login [post]
func PhoneLogin(c *gin.Context) {
	log.Info("Phone Login function called.")

	// Binding the data with the u struct.
	var req PhoneLoginCredentials
	if err := c.Bind(&req); err != nil {
		log.Warnf("phone login bind param err: %v", err)
		response.Error(c, errcode.ErrInvalidParam)
		return
	}

	log.Infof("req %#v", req)
	// check param
	if req.Phone == 0 || req.VerifyCode == 0 {
		log.Warn("phone login bind param is empty")
		response.Error(c, errcode.ErrInvalidParam)
		return
	}

	// 验证校验码
	if !service.Svc.VCode().CheckLoginVCode(req.Phone, req.VerifyCode) {
		response.Error(c, ecode.ErrVerifyCode)
		return
	}

	// 登录
	t, err := service.Svc.Users().PhoneLogin(c, req.Phone, req.VerifyCode)
	if err != nil {
		response.Error(c, ecode.ErrVerifyCode)
		return
	}

	response.Success(c, model.Token{Token: t})
}
