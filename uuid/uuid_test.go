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
		Run(func (_ *testing.T, data test) {
			u, err := FromString(data.arg)
			if err != nil {
				t.Errorf("From String: %v\n", err)
				return
			}

			got := u.String()
			y.AssertEqual(t, got, data.arg)
		})
}
