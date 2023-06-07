package codec

import (
	"go-im/common/proto/message"
	"io"
)

type Codec interface {
	io.Closer
	ReadData() (*message.Cmd, error)
	WriteData(cmd *message.Cmd) error
}
