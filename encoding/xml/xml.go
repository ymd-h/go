// Package encoding/xml implements Encoder/Decoder for XML
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

// Encode encodes data and returns encoded io.Reader.
func (_ Encoder) Encode(data any) (io.Reader, error) {
	return encoding.Encode(xml.NewEncoder, data)
}

// Decode decodes buf io.Reader to ptr.
func (_ Decoder) Decode(buf io.Reader, ptr any) error {
	return encoding.Decode(xml.NewDecoder, buf, ptr)
}
