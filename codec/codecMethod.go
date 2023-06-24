package codec

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
)

func init() {
	NewCodecFuncMap = make(map[Type]NewCodecFunc)
	NewCodecFuncMap[GobCodec] = NewGobCodec
	NewCodecFuncMap[JsonCodec] = NewJsonCodec
}

// 传递编解码器类型，返回对应的编解码器
func ParseOption(opt []byte) (NewCodecFunc, error) {
	if len(opt) != 5 {
		log.Println("codec error: error option len")
		return nil, errors.New("codec error: error option len")
	}
	optNum := binary.BigEndian.Uint32(opt[0:4])
	// log.Println("optNum: ", optNum)
	if optNum != magicNum {
		log.Println("codec error: magic number mismatch")
		return nil, errors.New("codec error: magic number mismatch")
	}
	newCodecFunc, ok := NewCodecFuncMap[Type(opt[4])]
	if !ok {
		log.Println("codec error: codec type not supported")
		return nil, errors.New("codec error: codec type not supported")
	}
	return newCodecFunc, nil
}

// 获取编解码器类型，用于传递给服务端
// 使用字节切片来传递编解码器类型
func GetOption(t Type) []byte {
	opt := make([]byte, 0)
	buf := bytes.NewBuffer(opt)
	_ = binary.Write(buf, binary.BigEndian, magicNum)
	_ = binary.Write(buf, binary.BigEndian, t)
	return buf.Bytes()
}
