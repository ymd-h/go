package xml

import (
	"bytes"

	"github.com/ymd-h/go/encoding/xml"
)

type (
	Encoder struct {}
	Decoder = xml.Decoder
)


func (_ Encoder) Encode(body any) (*bytes.Buffer, error) {
	return xml.Encoder().Encode(body)
}

func (_ Encoder) ContentType() string {
	return 'application/xml: charset="UTF-8"'
}
