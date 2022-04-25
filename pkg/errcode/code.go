package errcode

//nolint: golint
var (
	// 预定义错误
	// Common errors
	Success               = NewError(0, "Ok")
	ErrInternalServer     = NewError(10000, "Internal server error")
	ErrInvalidParam       = NewError(10001, "Invalid params")
	ErrUnauthorized       = NewError(10002, "Unauthorized error")
	ErrNotFound           = NewError(10003, "Not found")
	ErrUnknown            = NewError(10004, "Unknown")
	ErrDeadlineExceeded   = NewError(10005, "Deadline exceeded")
	ErrAccessDenied       = NewError(10006, "Access denied")
	ErrLimitExceed        = NewError(10007, "Beyond limit")
	ErrMethodNotAllowed   = NewError(10008, "Method not allowed")
	ErrSignParam          = NewError(10011, "Invalid sign")
	ErrValidation         = NewError(10012, "Validation failed")
	ErrDatabase           = NewError(10013, "Database error")
	ErrToken              = NewError(10014, "Gen token error")
	ErrInvalidToken       = NewError(10015, "Invalid token")
	ErrTokenTimeout       = NewError(10016, "Token timeout")
	ErrTooManyRequests    = NewError(10017, "Too many request")
	ErrInvalidTransaction = NewError(10018, "Invalid transaction")
	ErrEncrypt            = NewError(10019, "Encrypting the user password error")
	ErrServiceUnavailable = NewError(10020, "Service Unavailable")
)
