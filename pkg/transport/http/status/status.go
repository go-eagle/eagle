package status

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

const (
	// ClientClosed is non-standard http status code, which defined by nginx.
	// https://httpstatus.in/499/
	ClientClosed = 499
)

// Converter is a status converter.
type Converter interface {
	// GRPCCodeFromStatus converts a HTTP error code into the corresponding gRPC response status.
	GRPCCodeFromStatus(code int) codes.Code

	// HTTPStatusFromCode converts a gRPC error code into the corresponding HTTP response status.
	HTTPStatusFromCode(code codes.Code) int
}

type statusConverter struct{}

// DefaultConverter is the default status converter.
var DefaultConverter Converter = &statusConverter{}

// GRPCCodeFromStatus converts a HTTP error code into the corresponding gRPC response status.
// See: https://github.com/googleapis/googleapis/blob/master/google/rpc/code.proto
func (c statusConverter) GRPCCodeFromStatus(code int) codes.Code {
	switch code {
	case http.StatusOK:
		return codes.OK
	case http.StatusBadRequest:
		return codes.InvalidArgument
	case http.StatusUnauthorized:
		return codes.Unauthenticated
	case http.StatusForbidden:
		return codes.PermissionDenied
	case http.StatusNotFound:
		return codes.NotFound
	case http.StatusConflict:
		return codes.Aborted
	case http.StatusTooManyRequests:
		return codes.ResourceExhausted
	case http.StatusInternalServerError:
		return codes.Internal
	case http.StatusNotImplemented:
		return codes.Unimplemented
	case http.StatusServiceUnavailable:
		return codes.Unavailable
	case http.StatusGatewayTimeout:
		return codes.DeadlineExceeded
	case ClientClosed:
		return codes.Canceled
	}

	return codes.Unknown
}

// HTTPStatusFromCode converts a gRPC error code into the corresponding HTTP response status.
func (c statusConverter) HTTPStatusFromCode(code codes.Code) int {
	switch code {
	case codes.OK:
		return http.StatusOK
	case codes.Canceled:
		return ClientClosed
	case codes.Unknown:
		return http.StatusInternalServerError
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.FailedPrecondition:
		return http.StatusBadRequest
	case codes.Aborted:
		return http.StatusConflict
	case codes.OutOfRange:
		return http.StatusBadRequest
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.DataLoss:
		return http.StatusInternalServerError
	}

	return http.StatusOK
}

// GRPCCodeFromStatus converts a HTTP error code into the corresponding gRPC response status.
func GRPCCodeFromStatus(code int) codes.Code {
	return DefaultConverter.GRPCCodeFromStatus(code)
}

// HTTPStatusFromCode converts a gRPC error code into the corresponding HTTP response status.
func HTTPStatusFromCode(code codes.Code) int {
	return DefaultConverter.HTTPStatusFromCode(code)
}
