package ecode

import "github.com/1024casts/snake/pkg/errno"

var (
	// user errors
	ErrUserNotFound          = &errno.Errno{Code: 20101, Message: "The user was not found."}
	ErrPasswordIncorrect     = &errno.Errno{Code: 20102, Message: "账号或密码错误"}
	ErrAreaCodeEmpty         = &errno.Errno{Code: 20103, Message: "手机区号不能为空"}
	ErrPhoneEmpty            = &errno.Errno{Code: 20104, Message: "手机号不能为空"}
	ErrGenVCode              = &errno.Errno{Code: 20105, Message: "生成验证码错误"}
	ErrSendSMS               = &errno.Errno{Code: 20106, Message: "发送短信错误"}
	ErrSendSMSTooMany        = &errno.Errno{Code: 20107, Message: "已超出当日限制，请明天再试"}
	ErrVerifyCode            = &errno.Errno{Code: 20108, Message: "验证码错误"}
	ErrEmailOrPassword       = &errno.Errno{Code: 20109, Message: "邮箱或密码错误"}
	ErrTwicePasswordNotMatch = &errno.Errno{Code: 20110, Message: "两次密码输入不一致"}
	ErrRegisterFailed        = &errno.Errno{Code: 20111, Message: "注册失败"}
)
