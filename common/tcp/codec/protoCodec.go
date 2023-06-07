package codec

import (
	"bytes"
	"encoding/binary"
	"go-im/common/proto/message"
	"google.golang.org/protobuf/proto"
	"io"
)

type DataPkg struct {
	Len  int64 //8字节
	Data []byte
}

type ProtoCodec struct {
	conn io.ReadWriteCloser
}

func (p *ProtoCodec) Close() error {
	return p.conn.Close()
}

func NewProtoCodec(conn io.ReadWriteCloser) Codec {
	return &ProtoCodec{
		conn: conn,
	}
}

func (p *ProtoCodec) ReadData() (*message.Cmd, error) {

	// 读取4个字节的Len
	dataLenBuf := make([]byte, 4)
	if err := p.readFixedBytes(dataLenBuf); err != nil {
		return nil, err
	}
	var dataLen uint32
	if err := binary.Read(bytes.NewBuffer(dataLenBuf), binary.BigEndian, &dataLen); err != nil {
		return nil, err
	}
	// 读取data
	data := make([]byte, dataLen)
	if err := p.readFixedBytes(data); err != nil {
		return nil, err
	}
	cmd := &message.Cmd{}
	if err := proto.Unmarshal(data, cmd); err != nil {
		return nil, err
	}
	return cmd, nil
}

func (p *ProtoCodec) WriteData(cmd *message.Cmd) error {

	data, err := proto.Marshal(cmd)
	if err != nil {
		return err
	}

	//计算data长度，大端存储方式将data长度附加在data的前面
	dataBuf := bytes.NewBuffer([]byte{})
	if err := binary.Write(dataBuf, binary.BigEndian, uint32(len(data))); err != nil {
		return nil
	}
	dataBytes := append(dataBuf.Bytes(), data...)
	pos := 0
	for {
		n, err := p.conn.Write(dataBytes[pos:])
		if err != nil {
			return err
		}
		pos += n
		if pos == len(dataBytes) {
			return nil
		}
	}
}

// 读取固定字节
func (p *ProtoCodec) readFixedBytes(buf []byte) error {
	var pos int
	totalLen := len(buf)
	for {
		n, err := p.conn.Read(buf[pos:])
		if err != nil {
			return err
		}
		pos += n
		if pos == totalLen {
			break
		}
	}
	return nil
}
