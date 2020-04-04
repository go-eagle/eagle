package cache

import (
	"bytes"
	"compress/gzip"
	"encoding"
	"encoding/gob"

	"encoding/json"
	"io/ioutil"

	"github.com/golang/snappy"

	//json "github.com/json-iterator/go"
	"github.com/vmihailenco/msgpack"
)

type Encoding interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
}

func Marshal(e Encoding, v interface{}) (data []byte, err error) {
	bm, ok := v.(encoding.BinaryMarshaler)
	if ok && e == nil {
		data, err = bm.MarshalBinary()
		return
	} else {
		data, err = e.Marshal(v)
		if err == nil {
			return
		}
		if ok {
			data, err = bm.MarshalBinary()
		}
	}
	return
}

func Unmarshal(e Encoding, data []byte, v interface{}) (err error) {
	bm, ok := v.(encoding.BinaryUnmarshaler)
	if ok && e == nil {
		err = bm.UnmarshalBinary(data)
		return err
	} else {
		err = e.Unmarshal(data, v)
		if err == nil {
			return
		}
		if ok {
			return bm.UnmarshalBinary(data)
		}
	}
	return
}

type JsonEncoding struct{}

func (this JsonEncoding) Marshal(v interface{}) ([]byte, error) {
	buf, err := json.Marshal(v)
	return buf, err
}

func (this JsonEncoding) Unmarshal(data []byte, value interface{}) error {
	err := json.Unmarshal(data, value)
	if err != nil {
		return err
	}
	return nil
}

type GobEncoding struct{}

func (this GobEncoding) Marshal(v interface{}) ([]byte, error) {
	var (
		buffer bytes.Buffer
	)

	err := gob.NewEncoder(&buffer).Encode(v)
	return buffer.Bytes(), err
}

func (this GobEncoding) Unmarshal(data []byte, value interface{}) error {
	err := gob.NewDecoder(bytes.NewReader(data)).Decode(value)
	if err != nil {
		return err
	}
	return nil
}

type JsonGzipEncoding struct{}

func (this JsonGzipEncoding) Marshal(v interface{}) ([]byte, error) {
	buf, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	// var bufSizeBefore = len(buf)

	buf, err = GzipEncode(buf)
	// log.Infof("gzip_json_compress_ratio=%d/%d=%.2f", bufSizeBefore, len(buf), float64(bufSizeBefore)/float64(len(buf)))
	return buf, err
}

func (this JsonGzipEncoding) Unmarshal(data []byte, value interface{}) error {
	jsonData, err := GzipDecode(data)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonData, value)
	if err != nil {
		return err
	}
	return nil
}

func GzipEncode(in []byte) ([]byte, error) {
	var (
		buffer bytes.Buffer
		out    []byte
		err    error
	)
	writer, err := gzip.NewWriterLevel(&buffer, gzip.BestCompression)
	if err != nil {
		return nil, err
	}

	_, err = writer.Write(in)
	if err != nil {
		writer.Close()
		return out, err
	}
	err = writer.Close()
	if err != nil {
		return out, err
	}

	return buffer.Bytes(), nil
}

func GzipDecode(in []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(in))
	if err != nil {
		var out []byte
		return out, err
	}
	defer reader.Close()

	return ioutil.ReadAll(reader)
}

// JSONSnappyEncoding json格式和snappy压缩
type JSONSnappyEncoding struct{}

// Marshal 序列化
func (s JSONSnappyEncoding) Marshal(v interface{}) (data []byte, err error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	d := snappy.Encode(nil, b)
	return d, nil
}

// Unmarshal 反序列化
func (s JSONSnappyEncoding) Unmarshal(data []byte, value interface{}) error {
	b, err := snappy.Decode(nil, data)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, value)
}

type MsgPackEncoding struct{}

func (this MsgPackEncoding) Marshal(v interface{}) ([]byte, error) {
	buf, err := msgpack.Marshal(v)
	return buf, err
}

func (this MsgPackEncoding) Unmarshal(data []byte, value interface{}) error {
	err := msgpack.Unmarshal(data, value)
	if err != nil {
		return err
	}
	return nil
}
