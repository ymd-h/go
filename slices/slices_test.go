package slices

import (
	"reflect"
	"testing"

	y "github.com/ymd-h/go/testing"
)

func TestNewSlice(t *testing.T) {
	type test struct {
		arg any
		want any
	}
	y.NewTest[test](t).
		Add("int", test{
			arg: 1,
			want: NewComparableSlice[int](),
		}).
		Add("string", test{
			arg: "a",
			want: NewComparableSlice[string](),
		}).
		Add("uint", test{
			arg: uint(1),
			want: NewComparableSlice[uint](),
		}).
		Add("*int (unsupported)", test{
			arg: func() *int {
				i := 1
				return &i
			}(),
			want: &Slice[*int]{ item: []*int{} },
		}).
		Run(func (_ *testing.T, data test) {
			wantName := reflect.TypeOf(data.want).Name()
			var got string

			switch any(data.arg).(type) {
			case int:
				got = reflect.TypeOf(NewSlice[int]()).Name()
			case string:
				got = reflect.TypeOf(NewSlice[string]()).Name()
			case uint:
				got = reflect.TypeOf(NewSlice[uint]()).Name()
			case *int:
				got = reflect.TypeOf(NewSlice[*int]()).Name()
			}
			y.AssertEqual(t, got, wantName)
		})
}
