package errcode

import (
	"fmt"
	"net/http"
	"sync"
)

// Error 返回错误码和消息的结构体
// nolint: govet
type Error struct {
	code    int      `json:"code"`
	msg     string   `json:"msg"`
	details []string `json:"details"`
}

var errorCodes = map[int]struct{}{}
var toStatus sync.Map

// NewError create a error
func NewError(code int, msg string) *Error {
	if _, ok := errorCodes[code]; ok {
		panic(fmt.Sprintf("code %d is exsit, please change one", code))
	}
	errorCodes[code] = struct{}{}
	return &Error{code: code, msg: msg}
}

// Error return a error string
func (e Error) Error() string {
	return fmt.Sprintf("code: %d, msg: %s", e.Code(), e.Msg())
}

// Code return error code
func (e *Error) Code() int {
	return e.code
}

// Msg return error msg
func (e *Error) Msg() string {
	return e.msg
}

// Msgf format error string
func (e *Error) Msgf(args []interface{}) string {
	return fmt.Sprintf(e.msg, args...)
}

// Details return more error details
func (e *Error) Details() []string {
	return e.details
}

// WithDetails return err with detail
func (e *Error) WithDetails(details ...string) *Error {
	newError := *e
	newError.details = []string{}
	newError.details = append(newError.details, details...)

	return &newError
}

// SetHTTPStatusCode set a specific http status code to err
func SetHTTPStatusCode(err *Error, status int) {
	toStatus.Store(err.Code(), status)
}

// ToHTTPStatusCode convert custom error code to http status code and avoid return unknown status code.
func ToHTTPStatusCode(code int) int {
	if status, ok := toStatus.Load(code); ok {
		return status.(int)
	}

	return http.StatusBadRequest
}

// Err represents an error
type Err struct {
	Code    int
	Message string
	Err     error
}

// Error return error string
func (err *Err) Error() string {
	return fmt.Sprintf("Err - code: %d, message: %s, error: %s", err.Code, err.Message, err.Err)
}

// DecodeErr 对错误进行解码，返回错误code和错误提示
func DecodeErr(err error) (int, string) {
	if err == nil {
		return Success.code, Success.msg
	}

	switch typed := err.(type) {
	case *Err:
		return typed.Code, typed.Message
	case *Error:
		return typed.code, typed.msg
	default:
	}

	return ErrInternalServer.Code(), err.Error()
}

func initToStatus() {
	for code, status := range map[int]int{
		Success.Code():               http.StatusOK,
		ErrInternalServer.Code():     http.StatusInternalServerError,
		ErrNotFound.Code():           http.StatusNotFound,
		ErrInvalidParam.Code():       http.StatusBadRequest,
		ErrToken.Code():              http.StatusUnauthorized,
		ErrInvalidToken.Code():       http.StatusUnauthorized,
		ErrTokenTimeout.Code():       http.StatusUnauthorized,
		ErrTooManyRequests.Code():    http.StatusTooManyRequests,
		ErrServiceUnavailable.Code(): http.StatusServiceUnavailable,
	} {
		toStatus.Store(code, status)
	}
}

func init() {
	initToStatus()
}
