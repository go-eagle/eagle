package encoding

import (
	"github.com/bytedance/sonic"
	"github.com/golang/snappy"
)

type SonicEncoding struct {
}

func (s SonicEncoding) Marshal(v interface{}) ([]byte, error) {
	buf, err := sonic.Marshal(v)
	return buf, err
}

func (s SonicEncoding) Unmarshal(data []byte, value interface{}) error {
	err := sonic.Unmarshal(data, value)
	if err != nil {
		return err
	}
	return nil
}

// SonicSnappyEncoding  Sonic json格式和snappy压缩
type SonicSnappyEncoding struct{}

// Marshal 序列化
func (s SonicSnappyEncoding) Marshal(v interface{}) (data []byte, err error) {
	b, err := sonic.Marshal(v)
	if err != nil {
		return nil, err
	}
	d := snappy.Encode(nil, b)
	return d, nil
}

// Unmarshal 反序列化
func (s SonicSnappyEncoding) Unmarshal(data []byte, value interface{}) error {
	b, err := snappy.Decode(nil, data)
	if err != nil {
		return err
	}

	return sonic.Unmarshal(b, value)
}
