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
