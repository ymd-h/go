package gob

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
)

type (
	Encoder struct {}
	Decoder struct {}
)


func (_ Encoder) Encode(data any) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer([]byte{})

	if data != nil {
		err := gob.NewEncoder(buf).Encode(data)
		if err != nil {
			return nil, fmt.Errorf("Fail to Encode Gob: %w", err)
		}
	}

	return buf, nil
}


func (_ Decoder) Decode(buf io.Reader, ptr any) error {
	if ptr == nil {
		return nil
	}

	err := gob.NewDecoder(buf).Decode(ptr)
	if err != nil {
		return fmt.Errorf("Fail to Decode Gob: %w", err)
	}

	return nil
}

