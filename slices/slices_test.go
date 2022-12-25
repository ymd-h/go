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
			want: &OrderedSlice[int]{
				ComparableSlice[int]{
					Slice[int]{
						item: []int{},
					},
				},
			},
		}).
		Add("string", test{
			arg: "a",
			want: &OrderedSlice[string]{
				ComparableSlice[string]{
					Slice[string]{
						item: []string{},
					},
				},
			},
		}).
		Add("uint", test{
			arg: uint(1),
			want: &OrderedSlice[uint]{
				ComparableSlice[uint]{
					Slice[uint]{
						item: []uint{},
					},
				},
			},
		}).
		Add("bool", test{
			arg: true,
			want: &ComparableSlice[bool]{ Slice[bool]{ item: []bool{} } },
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


func TestGet(t *testing.T) {
	type test struct{
		init ISlice[string]
		arg int
		want string
	}
	y.NewTest[test](t).
		Add("simple", test{
			init: func() ISlice[string] {
				s := NewSlice[string]()
				s.Append("a", "b")
				return s
			}(),
			arg: 1,
			want: "b",
		}).
		Run(func(_ *testing.T, data test) {
			y.AssertEqual(t, data.init.Get(data.arg), data.want)
		})
}


func TestSet(t *testing.T) {
	type test struct{
		init ISlice[int]
		idx int
		v int
		before int
		after int
	}
	y.NewTest[test](t).
		Add("simple", test{
			init: NewSliceFrom[int]([]int{1, 2}),
			idx: 0,
			v: 3,
			before: 1,
			after: 3,
		}).
		Run(func(_ *testing.T, data test) {
			m := data.init
			idx := data.idx
			y.AssertEqual(t, m.Get(idx), data.before)
			m.Set(idx, data.v)
			y.AssertEqual(t, m.Get(idx), data.after)
		})
}


func TestAppend(t *testing.T) {
	type test struct{
		init ISlice[int]
		arg []int
		want []int
	}
	y.NewTest[test](t).
		Add("simple", test{
			init: NewSlice[int](),
			arg: []int{1, 2, 3},
			want: []int{1, 2, 3},
		}).
		Add("append", test{
			init: &ComparableSlice[int]{
				Slice[int]{ item: []int{1, 2, 3} },
			},
			arg: []int{4, 5, 6},
			want: []int{1, 2, 3, 4, 5, 6},
		}).
		Run(func(_ *testing.T, data test) {
			m := data.init
			m.Append(data.arg...)
			for i, e := range data.want {
				y.AssertEqual(t, m.Get(i), e)
			}
		})
}


func TestSize(t *testing.T) {
	type test struct {
		init ISlice[string]
		want int
	}
	y.NewTest[test](t).
		Add("empty", test{
			init: NewSlice[string](),
			want: 0,
		}).
		Add("some", test{
			init: NewSliceFrom([]string{"a"}),
			want: 1,
		}).
		Run(func(_ *testing.T, data test) {
			y.AssertEqual(t, data.init.Size(), data.want)
		})
}
