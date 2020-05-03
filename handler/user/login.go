package user

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/1024casts/snake/handler"
	"github.com/1024casts/snake/model"
	"github.com/1024casts/snake/pkg/errno"
	"github.com/1024casts/snake/pkg/log"
	"github.com/1024casts/snake/pkg/token"
	"github.com/1024casts/snake/service/user"
	"github.com/1024casts/snake/service/vcode"
)

// Login 手机登录接口
// @Summary 用户登录接口
// @Description 仅限手机登录
// @Tags 用户
// @Produce  json
// @Param req body PhoneLoginCredentials true "phone"
// @Success 200 {string} json "{"code":0,"message":"OK","data":{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6Ik"}}"
// @Router /users/login [post]
func Login(c *gin.Context) {
	log.Info("Phone Login function called.")

	// Binding the data with the u struct.
	var req PhoneLoginCredentials
	if err := c.Bind(&req); err != nil {
		handler.SendResponse(c, errno.ErrBind, nil)
		return
	}

	log.Infof("req %#v", req)
	// check param
	if req.Phone == 0 || req.VerifyCode == 0 {
		handler.SendResponse(c, errno.ErrParam, nil)
		return
	}

	// 验证校验码
	if !vcode.VCodeService.CheckLoginVCode(req.Phone, req.VerifyCode) {
		handler.SendResponse(c, errno.ErrVerifyCode, nil)
		return
	}

	// 如果是已经注册用户，则通过手机号获取用户信息
	u, err := user.UserSvc.GetUserByPhone(req.Phone)
	if err != nil {
		log.Warnf("[login] get u info err, %v", err)
	}

	// 否则新建用户信息, 并取得用户信息
	if u.ID == 0 {
		u := model.UserModel{
			Phone:    req.Phone,
			Username: strconv.Itoa(req.Phone),
		}
		u.ID, err = user.UserSvc.CreateUser(u)
		if err != nil {
			log.Warnf("[login] create u err, %v", err)
			handler.SendResponse(c, errno.InternalServerError, nil)
			return
		}
	}

	// 签发签名 Sign the json web token.
	t, err := token.Sign(c, token.Context{UserID: u.ID, Username: u.Username}, "")
	if err != nil {
		log.Warnf("[login] gen token sign err:, %v", err)
		handler.SendResponse(c, errno.ErrToken, nil)
		return
	}

	handler.SendResponse(c, nil, model.Token{
		Token: t,
	})
}
