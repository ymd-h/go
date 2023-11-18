// Package encoding implements unified interface of Encoder/Decoder.
package encoding

import (
	"bytes"
	"fmt"
	"io"
)

type (
	IEncoder interface {
		Encode(any) error
	}

	IDecoder interface {
		Decode(any) error
	}
)

// Encode encodes data with IEncoder and returns encoded io.Reader.
// If data is nil, nil io.Reader is returned.
// The encoding/json.NewEncoder in standard library can be passed.
func Encode[E IEncoder](newEncoder func(io.Writer) E, data any) (io.Reader, error) {
	if data == nil {
		return nil, nil
	}

	buf := bytes.NewBuffer([]byte{})

	err := newEncoder(buf).Encode(data)
	if err != nil {
		return nil, fmt.Errorf("Fail to Encode: %w", err)
	}

	return buf, nil
}

// Decode decodes buf io.Reader with IDecoder to ptr.
// The encoding/json.NewDecoder in standard library can be passed.
func Decode[D IDecoder](newDecoder func(io.Reader) D, buf io.Reader, ptr any) error {
	if ptr == nil {
		return nil
	}

	err := newDecoder(buf).Decode(ptr)
	if err != nil {
		return fmt.Errorf("Fail to Decode: %w", err)
	}

	return nil
}
