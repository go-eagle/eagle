package errno

//nolint: golint
var (
	// 预定义错误
	// Common errors
	OK                    = &Errno{Code: 0, Message: "OK"}
	InternalServerError   = &Errno{Code: 10001, Message: "Internal server error"}
	ErrBind               = &Errno{Code: 10002, Message: "Error occurred while binding the request body to the struct."}
	ErrParam              = &Errno{Code: 10003, Message: "Invalid params"}
	ErrSignParam          = &Errno{Code: 10004, Message: "Invalid sign"}
	ErrValidation         = &Errno{Code: 10005, Message: "Validation failed."}
	ErrDatabase           = &Errno{Code: 10006, Message: "Database error."}
	ErrToken              = &Errno{Code: 10007, Message: "Error occurred while signing the JSON web token."}
	ErrInvalidTransaction = &Errno{Code: 10008, Message: "Invalid transaction."}
	ErrTokenInvalid       = &Errno{Code: 10109, Message: "The token was invalid."}
	ErrEncrypt            = &Errno{Code: 10110, Message: "Error occurred while encrypting the user password."}
)
