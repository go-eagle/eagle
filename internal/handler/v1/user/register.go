package user

import (
	"github.com/1024casts/snake/internal/ecode"
	"github.com/1024casts/snake/internal/service"
	"github.com/gin-gonic/gin"

	"github.com/1024casts/snake/pkg/errno"
	"github.com/1024casts/snake/pkg/log"
)

// Register 注册
// @Summary 注册
// @Description 用户注册
// @Tags 用户
// @Produce  json
// @Param req body RegisterRequest true "请求参数"
// @Success 200 {object} model.UserInfo "用户信息"
// @Router /Register [post]
func Register(c *gin.Context) {
	// Binding the data with the u struct.
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warnf("register bind param err: %v", err)
		response.Error(c, errno.ErrBind)
		return
	}

	log.Infof("register req: %#v", req)
	// check param
	if req.Username == "" || req.Email == "" || req.Password == "" {
		log.Warnf("params is empty: %v", req)
		response.Error(c, errno.ErrInvalidParam)
		return
	}

	// 两次密码是否正确
	if req.Password != req.ConfirmPassword {
		log.Warnf("twice password is not same")
		response.Error(c, ecode.ErrTwicePasswordNotMatch)
		return
	}

	err := service.UserSvc.Register(c, req.Username, req.Email, req.Password)
	if err != nil {
		log.Warnf("register err: %v", err)
		response.Error(c, ecode.ErrRegisterFailed)
		return
	}

	response.Success(c, nil)
}
