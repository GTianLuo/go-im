package codec

import (
	"go-im/common/tcp"
	"io"
)

type Codec interface {
	io.Closer
	ReadFixedHeader(*tcp.FixedHeader) error
	ReadBody(interface{}) error
	Write(*tcp.FixedHeader, interface{}) error
}
