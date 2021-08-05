package errcode

//nolint: golint
var (
	// 预定义错误
	// Common errors
	Success               = NewError(0, "Success")
	ErrInternalServer     = NewError(10001, "Internal server error")
	ErrBind               = NewError(10002, "Bind request error")
	ErrInvalidParam       = NewError(10003, "Invalid params")
	ErrSignParam          = NewError(10004, "Invalid sign")
	ErrValidation         = NewError(10005, "Validation failed")
	ErrDatabase           = NewError(10006, "Database error")
	ErrToken              = NewError(10007, "Gen token error")
	ErrInvalidToken       = NewError(10108, "Invalid token")
	ErrTokenTimeout       = NewError(10109, "Token timeout")
	ErrTooManyRequests    = NewError(10110, "Too many request")
	ErrInvalidTransaction = NewError(10111, "Invalid transaction")
	ErrEncrypt            = NewError(10112, "Encrypting the user password error")
	ErrLimitExceed        = NewError(10113, "Beyond limit")
	ErrServiceUnavailable = NewError(10114, "Service Unavailable")
)
