package gob

import (
	"bytes"

	"github.com/ymd-h/go/encoding/gob"
)

type (
	Encoder struct {}
	Decoder = gob.Decoder
)


func (_ Encoder) Encode(body any) (*bytes.Buffer, error) {
	return gob.Encoder().Encode(body)
}

func (_ Encoder) ContentType() string {
	return "application/octet-stream"
}
