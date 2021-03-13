package ecode

import "github.com/1024casts/snake/pkg/errno"

var (
	// user errors
	ErrUserNotFound          = errno.NewError(20101, "The user was not found.")
	ErrPasswordIncorrect     = errno.NewError(20102, "账号或密码错误")
	ErrAreaCodeEmpty         = errno.NewError(20103, "手机区号不能为空")
	ErrPhoneEmpty            = errno.NewError(20104, "手机号不能为空")
	ErrGenVCode              = errno.NewError(20105, "生成验证码错误")
	ErrSendSMS               = errno.NewError(20106, "发送短信错误")
	ErrSendSMSTooMany        = errno.NewError(20107, "已超出当日限制，请明天再试")
	ErrVerifyCode            = errno.NewError(20108, "验证码错误")
	ErrEmailOrPassword       = errno.NewError(20109, "邮箱或密码错误")
	ErrTwicePasswordNotMatch = errno.NewError(20110, "两次密码输入不一致")
	ErrRegisterFailed        = errno.NewError(20111, "注册失败")
)
