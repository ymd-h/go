package xml

import (
	"bytes"
	"encoding/xml"
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
		err := xml.NewEncoder(buf).Encode(data)
		if err != nil {
			return nil, fmt.Errorf("Fail to Encode XML: %w", err)
		}
	}

	return buf, nil
}


func (_ Decoder) Decode(buf io.Reader, ptr any) error {
	if ptr == nil {
		return nil
	}

	err := xml.NewDecoder(buf).Decode(ptr)
	if err != nil {
		return fmt.Errorf("Fail to Decode XML: %w", err)
	}

	return nil
}
