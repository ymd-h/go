package xml

import (
	"encoding/xml"
	"io"

	"github.com/ymd-h/go/encoding"
)

type (
	Encoder struct {}
	Decoder struct {}
)


func (_ Encoder) Encode(data any) (io.Reader, error) {
	return encoding.Encode(xml.NewEncoder, data)
}


func (_ Decoder) Decode(buf io.Reader, ptr any) error {
	return encoding.Decode(xml.NewDecoder, buf, ptr)
}
