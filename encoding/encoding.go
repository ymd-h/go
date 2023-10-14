package encoding

import (
	"bytes"
	"io"
)

type (
	Encoder interface {
		Encode(any) (*bytes.Buffer, error)
	}

	Decoder interface {
		Decode(io.Reader, any) error
	}
)
