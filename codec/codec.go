package codec

import (
	"io"
)

type Header struct {
	ServiceMethod string // rpc格式为 "Service.Method"
	Seq           uint64 // 客户端请求的序号
	Error         string
}

type Codec interface {
	io.Closer                 // 关闭连接
	ReadHeader(*Header) error //读取头部
	ReadBody(interface{}) error
	Write(*Header, interface{}) error
}

// supporting codec
type Type uint8
type NewCodecFunc func(io.ReadWriteCloser) Codec

const (
	GobCodec Type = iota
	JsonCodec
	// 使用iota自增枚举值用于标识不同的编解码器
)

var NewCodecFuncMap map[Type]NewCodecFunc

// 用于标识RPC请求的魔数
const magicNum uint32 = 0x3a2fe23
