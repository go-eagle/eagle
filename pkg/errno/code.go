package errno

//nolint: golint
var (
	// 预定义错误
	// Common errors
	OK                  = &Errno{Code: 0, Message: "OK"}
	InternalServerError = &Errno{Code: 10001, Message: "Internal server error"}
	ErrBind             = &Errno{Code: 10002, Message: "Error occurred while binding the request body to the struct."}
	ErrParam            = &Errno{Code: 10003, Message: "参数有误"}
	ErrSignParam        = &Errno{Code: 10004, Message: "签名参数有误"}

	ErrValidation         = &Errno{Code: 20001, Message: "Validation failed."}
	ErrDatabase           = &Errno{Code: 20002, Message: "Database error."}
	ErrToken              = &Errno{Code: 20003, Message: "Error occurred while signing the JSON web token."}
	ErrInvalidTransaction = &Errno{Code: 20004, Message: "invalid transaction."}

	// user errors
	ErrEncrypt               = &Errno{Code: 20101, Message: "Error occurred while encrypting the user password."}
	ErrUserNotFound          = &Errno{Code: 20102, Message: "The user was not found."}
	ErrTokenInvalid          = &Errno{Code: 20103, Message: "The token was invalid."}
	ErrPasswordIncorrect     = &Errno{Code: 20104, Message: "账号或密码错误"}
	ErrAreaCodeEmpty         = &Errno{Code: 20105, Message: "手机区号不能为空"}
	ErrPhoneEmpty            = &Errno{Code: 20106, Message: "手机号不能为空"}
	ErrGenVCode              = &Errno{Code: 20107, Message: "生成验证码错误"}
	ErrSendSMS               = &Errno{Code: 20108, Message: "发送短信错误"}
	ErrSendSMSTooMany        = &Errno{Code: 20109, Message: "已超出当日限制，请明天再试"}
	ErrVerifyCode            = &Errno{Code: 20110, Message: "验证码错误"}
	ErrEmailOrPassword       = &Errno{Code: 20111, Message: "邮箱或密码错误"}
	ErrTwicePasswordNotMatch = &Errno{Code: 20112, Message: "两次密码输入不一致"}
	ErrRegisterFailed        = &Errno{Code: 20113, Message: "注册失败"}
)
