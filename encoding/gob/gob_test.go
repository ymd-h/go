package gob

import (
	"testing"

	"github.com/ymd-h/go/slices"
)

func TestJSON(t *testing.T){
	type (
		B struct {
			B1 bool
			B2 []uint16
		}
		A struct {
			A1 string
			A2 int
			A3 B
		}
	)

	a := A{
		A1: "12345abcde",
		A2: 255,
		A3: B{
			B1: true,
			B2: []uint16{0, 1, 2, 16},
		},
	}

	enc := Encoder{}
	dec := Decoder{}

	b, err := enc.Encode(a)
	if err != nil {
		t.Errorf("Fail: %v\n", err)
		return
	}

	var aa A
	err = dec.Decode(b, &aa)
	if err != nil {
		t.Errorf("Fail: %v\n", err)
		return
	}

	if (a.A1 != aa.A1) ||
		(a.A2 != aa.A2) ||
		(a.A3.B1 != aa.A3.B1) ||
		(!slices.NewComparableSliceFrom(a.A3.B2).Equal(
			slices.NewComparableSliceFrom(aa.A3.B2),
		)) {
		t.Errorf("Fail\n")
		return
	}

}
