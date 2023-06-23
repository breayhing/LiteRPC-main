package codec

import (
	"encoding/json"
	"io"
	"log"
)

type jsonCodec struct {
	conn io.ReadWriteCloser
	enc  *json.Encoder
	dec  *json.Decoder
}

func NewJsonCodec(conn io.ReadWriteCloser) Codec {
	return &jsonCodec{
		conn: conn,
		enc:  json.NewEncoder(conn),
		dec:  json.NewDecoder(conn),
	}
}

func (c *jsonCodec) ReadHeader(header *Header) error {
	return c.dec.Decode(header)
}

func (c *jsonCodec) ReadBody(body interface{}) error {
	return c.dec.Decode(body)
}

func (c *jsonCodec) Write(header *Header, body interface{}) (err error) {
	defer func() {
		if err != nil {
			_ = c.Close()
		}
	}()
	if err := c.enc.Encode(header); err != nil {
		log.Println("codec error: json error encoding header:", err)
		return err
	}
	if err := c.enc.Encode(body); err != nil {
		log.Println("codec error: json error encoding body:", err)
		return err
	}
	return nil
}

func (c *jsonCodec) Close() error {
	return c.conn.Close()
}
