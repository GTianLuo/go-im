package codec

import (
	"bufio"
	"encoding/gob"
	"go-im/common/tcp"
	"io"
)

type GobCodec struct {
	conn io.ReadWriteCloser
	buf  *bufio.Writer //写缓冲
	dec  *gob.Decoder
	enc  *gob.Encoder
}

func NewGobCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	return &GobCodec{
		conn: conn,
		buf:  buf,
		dec:  gob.NewDecoder(conn),
		enc:  gob.NewEncoder(buf),
	}
}

func (g *GobCodec) Close() error {
	return g.conn.Close()
}

func (g *GobCodec) ReadFixedHeader(header *tcp.FixedHeader) error {

	return g.dec.Decode(header)
}

func (g *GobCodec) ReadBody(i interface{}) error {
	return g.dec.Decode(i)
}

func (g *GobCodec) Write(header *tcp.FixedHeader, i interface{}) error {
	defer func() { _ = g.buf.Flush() }()
	if err := g.enc.Encode(header); err != nil {
		return err
	}
	if err := g.enc.Encode(i); err != nil {
		return err
	}
	return nil
}
