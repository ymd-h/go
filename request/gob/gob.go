package gob

import (
	"github.com/ymd-h/go/encoding/gob"
)

type (
	Encoder struct {
		gob.Encoder
	}
	Decoder = gob.Decoder
)


func (_ Encoder) ContentType() string {
	return "application/octet-stream"
}
