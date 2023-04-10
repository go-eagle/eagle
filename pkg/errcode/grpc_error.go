package errcode

import (
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

// GrpcStatus grpc error
type GrpcStatus struct {
	status  *status.Status
	details []proto.Message
}

// New instance a status
func New(code codes.Code, msg string) *GrpcStatus {
	return &GrpcStatus{
		status: status.New(code, msg),
	}
}

// Status .
func (g *GrpcStatus) Status(details ...proto.Message) *status.Status {
	details = append(details, g.details...)
	st, err := g.status.WithDetails(details...)
	if err != nil {
		return g.status
	}
	return st
}

// WithDetails .
func (g *GrpcStatus) WithDetails(details ...proto.Message) *GrpcStatus {
	g.details = details
	return g
}

// NewDetails .
func NewDetails(details map[string]interface{}) proto.Message {
	detailStruct, err := structpb.NewStruct(details)
	if err != nil {
		return nil
	}
	return detailStruct
}

// ToRPCCode 自定义错误码转换为RPC识别的错误码，避免返回Unknown状态码
func ToRPCCode(code int) codes.Code {
	var statusCode codes.Code

	switch code {
	case ErrInternalServer.code:
		statusCode = codes.Internal
	case ErrInvalidParam.code:
		statusCode = codes.InvalidArgument
	case ErrUnauthorized.code:
		statusCode = codes.Unauthenticated
	case ErrNotFound.code:
		statusCode = codes.NotFound
	case ErrDeadlineExceeded.code:
		statusCode = codes.DeadlineExceeded
	case ErrAccessDenied.code:
		statusCode = codes.PermissionDenied
	case ErrLimitExceed.code:
		statusCode = codes.ResourceExhausted
	case ErrMethodNotAllowed.code:
		statusCode = codes.Unimplemented
	default:
		statusCode = codes.Unknown
	}

	return statusCode
}
