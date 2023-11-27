package uuid

import (
	"fmt"
	"testing"

	y "github.com/ymd-h/go/testing"
)

func TestText(t *testing.T) {
	type test struct {
		arg string
		ok bool
	}

	y.NewTest[test](t).
		Add("zero", test{
			arg: "00000000-0000-0000-0000-000000000000",
			ok: true,
		}).
		Add("short", test{
			arg: "00000000-0000-0000-0000-00000000000",
			ok: false,
		}).
		Add("long", test{
			arg: "00000000-0000-0000-0000-0000000000000",
			ok: false,
		}).
		Add("invalid", test{
			arg: "00000000-0000-x000-0000-000000000000",
			ok: false,
		}).
		Add("all", test{
			arg: "01234567-89ab-cdef-0123-456789abcdef",
			ok: true,
		}).
		Run(func (_ *testing.T, data test) {
			u, err := FromString(data.arg)
			if data.ok {
				if err != nil {
					t.Errorf("Fail: From String: %v\n", err)
					return
				}
			} else {
				if err == nil {
					t.Errorf("Must Faill: From String\n")
				}
				return
			}

			got := u.String()
			y.AssertEqual(t, got, data.arg)
		})
}

func testVersion[U interface {
	Version() uint8
	Variant() uint8
}](
	t *testing.T,
	f func() (U, error),
	version uint8,
	variant uint8,
) error {
	u, err := f()
	if err != nil {
		t.Errorf("Fail: %v\n", err)
		return err
	}

	if u.Version() != version {
		t.Errorf("Fail Version: %v != %v\n", u.Version(), version)
		return fmt.Errorf("Version Error")
	}

	if u.Variant() != variant {
		t.Errorf("Fail Variant: %v != %v\n", u.Variant(), variant)
		return fmt.Errorf("Variant Error")
	}

	return nil
}

func TestVersion(t *testing.T){
	if testVersion(t, NewUUIDv4, 4, 0b10) != nil {
		return
	}
	if testVersion(t, NewUUIDv7, 7, 0b10) != nil {
		return
	}
}


func TestBinary(t *testing.T){
	tests := []struct{
		name string
		data []byte
		ok bool
	}{
		{
			name: "zero",
			data: []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},
			ok: true,
		},
		{
			name: "short",
			data: []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},
			ok: false,
		},
		{
			name: "long",
			data: []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},
			ok: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func (t *testing.T){
			var u UUID

			err := u.UnmarshalBinary(test.data)
			if err != nil {
				if test.ok {
					t.Errorf("Fail: %v\n", err)
				}
				return
			} else {
				if !test.ok {
					t.Errorf("Must Fail\n")
					return
				}
			}

			b, err := u.MarshalBinary()
			if err != nil {
				t.Errorf("Fail: %v\n", err)
				return
			}

			if len(b) != len(test.data) {
				t.Errorf("Fail\n")
				return
			}
			for i, d := range test.data {
				if b[i] != d {
					t.Errorf("Fail\n")
					return
				}
			}
		})
	}
}
