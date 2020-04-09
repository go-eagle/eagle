package user

import (
	. "github.com/1024casts/snake/handler"
	"github.com/1024casts/snake/pkg/errno"
	"github.com/1024casts/snake/service/sms"
	"github.com/1024casts/snake/service/vcode"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/pkg/errors"
)

// VCode 获取验证码
// @Summary 根据手机号获取校验码
// @Description Get an user by username
// @Tags 用户
// @Accept  json
// @Produce  json
// @Param area_code query string true "区域码，比如86"
// @Param phone query string true "手机号"
// @Success 200 {object} handler.Response
// @Router /users/vcode [get]
func VCode(c *gin.Context) {
	log.Info("VCode function called.")

	// 验证区号和手机号是否为空
	if c.Query("area_code") == "" {
		SendResponse(c, errno.ErrAreaCodeEmpty, nil)
		return
	}

	phone := c.Query("phone")
	if phone == "" {
		SendResponse(c, errno.ErrPhoneEmpty, nil)
		return
	}

	// TODO: 频率控制，以防攻击

	// 生成短信验证码
	verifyCode, err := vcode.VCodeService.GenLoginVCode(phone)
	if err != nil {
		log.Warnf("gen login verify code err, %v", errors.WithStack(err))
		SendResponse(c, errno.ErrGenVCode, nil)
		return
	}

	// 发送短信
	err = sms.ServiceSms.Send(phone, verifyCode)
	if err != nil {
		log.Warnf("send phone sms err, %v", errors.WithStack(err))
		SendResponse(c, errno.ErrSendSMS, nil)
		return
	}

	SendResponse(c, nil, nil)
}
