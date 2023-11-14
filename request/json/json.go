package json

import (
	"bytes"

	"github.com/ymd-h/go/encoding/json"
)

type (
	Encoder struct {}
	Decoder = json.Decoder
)


func (_ Encoder) Encode(body any) (*bytes.Buffer, error) {
	return json.Encoder{}.Encode(body)
}

func (_ Encoder) ContentType() string {
	return "application/json"
}
