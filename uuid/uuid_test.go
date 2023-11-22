package uuid

import (
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
