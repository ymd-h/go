package json

import (
	"encoding/json"
	"io"

	"github.com/ymd-h/go/encoding"
)

type (
	Encoder struct {}
	Decoder struct {}
)


func (_ Encoder) Encode(data any) (io.Reader, error) {
	return encoding.Encode(json.NewEncoder, data)
}


func (_ Decoder) Decode(buf io.Reader, ptr any) error {
	return encoding.Decode(json.NewDecoder, buf, ptr)
}
