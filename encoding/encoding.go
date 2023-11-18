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


func Encode[E IEncoder](newEncoder func(io.Writer) E, data any) (io.Reader, error) {
	if data == nil {
		return nil, nil
	}

	buf := bytes.NewBuffer([]byte{})

	err := newEncoder(buf).Encode(data)
	if err != nil {
		return nil, fmt.Errorf("Fail to Encode: %w\n", err)
	}

	return buf, nil
}


func Decode[D IDecoder](newDecoder func(io.Reader) D, buf io.Reader, ptr any) error {
	if ptr == nil {
		return nil
	}

	err := newDecoder(buf).Decode(ptr)
	if err != nil {
		return fmt.Errorf("Fail to Decode: %w\n", err)
	}

	return nil
}
