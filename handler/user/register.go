package user

import (
	"github.com/gin-gonic/gin"

	"github.com/1024casts/snake/handler"
	"github.com/1024casts/snake/pkg/errno"
	"github.com/1024casts/snake/pkg/log"
	"github.com/1024casts/snake/service/user"
)

// Register 注册
// @Summary 注册
// @Description 用户注册
// @Tags 用户
// @Produce  json
// @Param req body RegisterRequest true ""
// @Success 200 {string} json "{"code":0,"message":"OK","data":{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6Ik"}}"
// @Router /Register [post]
func Register(c *gin.Context) {
	// Binding the data with the u struct.
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warnf("register bind param err: %v", err)
		handler.SendResponse(c, errno.ErrBind, nil)
		return
	}

	log.Infof("register req: %#v", req)
	// check param
	if req.Username == "" || req.Email == "" || req.Password == "" {
		log.Warnf("params is empty: %v", req)
		handler.SendResponse(c, errno.ErrParam, nil)
		return
	}

	// 两次密码是否正确
	if req.Password != req.ConfirmPassword {
		log.Warnf("twice password is not same")
		handler.SendResponse(c, errno.ErrTwicePasswordNotMatch, nil)
		return
	}

	err := user.Svc.Register(c, req.Username, req.Email, req.Password)
	if err != nil {
		log.Warnf("register err: %v", err)
		handler.SendResponse(c, errno.ErrRegisterFailed, nil)
		return
	}

	handler.SendResponse(c, nil, nil)
}
