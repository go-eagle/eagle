package grpc

import (
	"github.com/golang/protobuf/proto"
)

type codec struct{}

func (c *codec) Marshal(v interface{}) ([]byte, error) {
	b, err := proto.Marshal(v.(proto.Message))
	sentBytes.Add(float64(len(b)))
	return b, err
}

func (c *codec) Unmarshal(data []byte, v interface{}) error {
	receivedBytes.Add(float64(len(data)))
	return proto.Unmarshal(data, v.(proto.Message))
}

func (c *codec) String() string {
	return "proto"
}
