package uuid

import (
	"testing"

	y "github.com/ymd-h/go/testing"
)

func TestH2B(t *testing.T) {
	type test struct {
		arg byte
		want byte
		ok bool
	}

	y.NewTest[test](t).
		Add("0", test{
			arg: 0x30,
			want: 0,
			ok: true,
		}).
		Add("invalid", test{
			arg: 0,
			want: 0,
			ok: false,
		}).
		Add("f", test{
			arg: []byte("f")[0],
			want: 0xf,
			ok: true,

		}).
		Add("C", test{
			arg: []byte("C")[0],
			want: 0xc,
			ok: true,
		}).
		Run(func(_ *testing.T, data test){
			b, err := h2b(data.arg)
			if data.ok {
				if err != nil {
					t.Errorf("Fail h2b: %v\n", err)
					return
				}
			} else {
				if err == nil {
					t.Errorf("Must Fail h2b\n")
				}
				return
			}
			y.AssertEqual(t, b, data.want)
		})
}

func TestHex(t *testing.T) {
	type test struct {
		arg []byte
		want byte
		ok bool
	}

	y.NewTest[test](t).
		Add("zero", test{
			arg: []byte{0x30, 0x30},
			want: 0,
			ok: true,
		}).
		Run(func (_ *testing.T, data test){
			b, err := decodeHEX(data.arg)
			if data.ok {
				if err != nil {
					t.Errorf("Fail to decode HEX: %v\n", err)
					return
				}
			} else {
				if err == nil {
					t.Errorf("Must fail decode HEX\n")
				}
				return
			}
			y.AssertEqual(t, b, data.want)
		})
}
