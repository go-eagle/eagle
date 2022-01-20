package errcode

import (
	"github.com/golang/protobuf/proto"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

type GrpcStatus struct {
	status  *status.Status
	details []proto.Message
}

func New(code grpcCodes.Code, msg string) *GrpcStatus {
	return &GrpcStatus{
		status: status.New(code, msg),
	}
}

func (g *GrpcStatus) Status(details ...proto.Message) *status.Status {
	details = append(details, g.details...)
	st, err := g.status.WithDetails(details...)
	if err != nil {
		return g.status
	}
	return st
}

func (g *GrpcStatus) WithDetails(details ...proto.Message) *GrpcStatus {
	g.details = details
	return g
}

func NewDetails(details map[string]interface{}) proto.Message {
	detailStruct, err := structpb.NewStruct(details)
	if err != nil {
		return nil
	}
	return detailStruct
}
