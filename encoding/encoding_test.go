package encoding

import (
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"io"
	"testing"
)


func TestEncoding(t *testing.T){
	type (
		A struct {
			A1 string `json:"a1" xml:"a1"`
			A2 uint64 `json:"a2" xml:"a2"`
		}
	)

	tests := []struct {
		name string
		enc func(A) (io.Reader, error)
		dec func(io.Reader, any) error
		a A
	}{
		{
			name: "JSON",
			enc: func(a A) (io.Reader, error) {
				return Encode(json.NewEncoder, a)
			},
			dec: func(r io.Reader, ptr any) error {
				return Decode(json.NewDecoder, r, ptr)
			},
			a: A{ A1: "aaaa", A2: 12345 },
		},
		{
			name: "XML",
			enc: func(a A) (io.Reader, error) {
				return Encode(xml.NewEncoder, a)
			},
			dec: func(r io.Reader, ptr any) error {
				return Decode(xml.NewDecoder, r, ptr)
			},
			a: A{ A1: "abcdef", A2: 980876 },
		},
		{
			name: "Gob",
			enc: func(a A) (io.Reader, error) {
				return Encode(gob.NewEncoder, a)
			},
			dec: func(r io.Reader, ptr any) error {
				return Decode(gob.NewDecoder, r, ptr)
			},
			a: A{ A1: "--x-c-909", A2: 78998765 },
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T){
			b, err := test.enc(test.a)
			if err != nil {
				t.Errorf("Fail: %v\n", err)
				return
			}

			var aa A
			err = test.dec(b, &aa)
			if err != nil {
				t.Errorf("Fail: %v\n", err)
				return
			}

			if test.a.A1 != aa.A1 {
				t.Errorf("Fail: %v != %v\n", test.a.A1, aa.A1)
				return
			}
			if test.a.A2 != aa.A2 {
				t.Errorf("Fail: %v != %v\n", test.a.A2, aa.A2)
				return
			}
		})
	}
}
