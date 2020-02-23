package user

import (
	"strconv"

	. "github.com/1024casts/snake/handler"
	"github.com/1024casts/snake/model"
	"github.com/1024casts/snake/pkg/errno"
	"github.com/1024casts/snake/pkg/token"
	"github.com/1024casts/snake/service"
	"github.com/lexkong/log"

	"github.com/gin-gonic/gin"
)

// @Summary 用户登录接口 Done
// @Description 仅限手机登录
// @Tags 用户
// @Produce  json
// @Param phone body string true "phone"
// @Param verifyCode body string true "verify_code"
// @Success 200 {string} json "{"code":0,"message":"OK","data":{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1MjgwMTY5MjIsImlkIjowLCJuYmYiOjE1MjgwMTY5MjIsInVzZXJuYW1lIjoiYWRtaW4ifQ.LjxrK9DuAwAzUD8-9v43NzWBN7HXsSLfebw92DKd1JQ"}}"
// @Router /users/login [post]
func Login(c *gin.Context) {
	log.Info("Phone Login function called.")

	// Binding the data with the user struct.
	var req PhoneLoginCredentials
	if err := c.Bind(&req); err != nil {
		SendResponse(c, errno.ErrBind, nil)
		return
	}

	log.Infof("req %#v", req)
	// check param
	if req.Phone == 0 || req.VerifyCode == 0 {
		SendResponse(c, errno.ErrParam, nil)
		return
	}

	// 验证校验码
	if !service.VCodeService.CheckLoginVCode(req.Phone, req.VerifyCode) {
		SendResponse(c, errno.ErrVerifyCode, nil)
		return
	}

	// 如果是已经注册用户，则通过手机号获取用户信息
	user, err := service.UserService.GetUserByPhone(req.Phone)
	if err != nil {
		log.Warnf("[login] get user info err, %v", err)
	}

	// 否则新建用户信息, 并取得用户信息
	if user.Id == 0 {
		u := model.UserModel{
			Phone:    req.Phone,
			Username: strconv.Itoa(req.Phone),
		}
		user.Id, err = service.UserService.CreateUser(u)
		if err != nil {
			log.Warnf("[login] create user err, %v", err)
			SendResponse(c, errno.InternalServerError, nil)
			return
		}
	}

	// 签发签名 Sign the json web token.
	t, err := token.Sign(c, token.Context{UserID: user.Id, Username: user.Username}, "")
	if err != nil {
		log.Warnf("[login] gen token sign err:, %v", err)
		SendResponse(c, errno.ErrToken, nil)
		return
	}

	SendResponse(c, nil, model.Token{
		Token: t,
	})
}
