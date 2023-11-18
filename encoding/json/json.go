// Package encoding/json implements Encoder/Decoder for JSON
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

// Encode encodes data and returns encoded io.Reader.
func (_ Encoder) Encode(data any) (io.Reader, error) {
	return encoding.Encode(json.NewEncoder, data)
}

// Decode decodes buf io.Reader to ptr.
func (_ Decoder) Decode(buf io.Reader, ptr any) error {
	return encoding.Decode(json.NewDecoder, buf, ptr)
}
