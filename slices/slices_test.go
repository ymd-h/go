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

func TestBinarySearch(t *testing.T) {
	type test struct {
		init IOrderedSlice[int]
		target int
		found bool
		idx int
	}
	y.NewTest[test](t).
		Add("simple", test{
			init: NewOrderedSliceFrom([]int{1, 2, 3, 4}),
			target: 3,
			found: true,
			idx: 2,
		}).
		Add("not found", test{
			init: NewOrderedSliceFrom([]int{1, 2, 3}),
			target: 4,
			found: false,
			idx: -1,
		}).
		Run(func(_ *testing.T, data test) {
			idx, ok := data.init.BinarySearch(data.target)
			y.AssertEqual(t, ok, data.found)
			if ok {
				y.AssertEqual(t, idx, data.idx)
			}
		})
}

func TestBinarySearchFunc(t *testing.T) {
	type test struct {
		init ISlice[int]
		target int
		found bool
		idx int
	}
	y.NewTest[test](t).
		Add("simple", test{
			init: NewSliceFrom([]int{1, 2, 3, 4}),
			target: 3,
			found: true,
			idx: 2,
		}).
		Add("not found", test{
			init: NewSliceFrom([]int{1, 2, 3, 4}),
			target: 5,
			found: false,
			idx: -1,
		}).
		Run(func(_ *testing.T, data test) {
			f := func(a, b int) int {
				switch {
				case a > b:
					return 1
				case a < b:
					return -1
				default:
					return 0
				}
			}
			idx, ok := data.init.BinarySearchFunc(data.target, f)
			y.AssertEqual(t, ok, data.found)
			if ok {
				y.AssertEqual(t, idx, data.idx)
			}
		})
}


func TestClip(t *testing.T) {}

func TestClone(t *testing.T) {
	type test struct {
		init ISlice[int]
	}
	y.NewTest[test](t).
		Add("some", test{
			init: NewSliceFrom([]int{1, 2, 3}),
		}).
		Add("empty", test{
			init: NewSlice[int](),
		}).
		Run(func(_ *testing.T, data test) {
			c := data.init.Clone()
			y.AssertEqual(t, data.init.Size(), c.Size())
			for i:= 0; i < c.Size(); i++ {
				y.AssertEqual(t, data.init.Get(i), c.Get(i))
			}
		})
}


func TestCompact(t *testing.T) {
	type test struct {
		init IComparableSlice[int]
		want IComparableSlice[int]
	}
	y.NewTest[test](t).
		Add("simple", test{
			init: NewComparableSliceFrom([]int{1, 1, 2, 3, 4, 4}),
			want: NewComparableSliceFrom([]int{1, 2, 3, 4}),
		}).
		Add("same", test{
			init: NewComparableSliceFrom([]int{1, 2, 3}),
			want: NewComparableSliceFrom([]int{1, 2, 3}),
		}).
		Add("empty", test{
			init: NewComparableSlice[int](),
			want: NewComparableSlice[int](),
		}).
		Add("not sort", test{
			init: NewComparableSliceFrom([]int{1, 3, 2, 3, 2}),
			want: NewComparableSliceFrom([]int{1, 3, 2, 3, 2}),
		}).
		Run(func(_ *testing.T, data test) {
			data.init.Compact()
			y.AssertEqual(t, data.init.Equal(data.want), true)
		})
}


func TestCompactFunc(t *testing.T) {}

func TestContains(t *testing.T) {}

func TestContainsFunc(t *testing.T) {}

func TestDelete(t *testing.T) {}

func TestEqual(t *testing.T) {}

func TestEqualFunc(t *testing.T) {}

func TestGrow(t *testing.T) {}

func TestIndex(t *testing.T) {}

func TestIndexFunc(t *testing.T) {}

func TestInsert(t *testing.T) {}

func TestIsSorted(t *testing.T) {}

func TestIsSortedFunc(t *testing.T) {}

func TestReplace(t *testing.T) {}

func TestSort(t *testing.T) {}

func TestSortFunc(t *testing.T) {}

func TestTryComparable(t *testing.T) {}

func TestTryOrdered(t *testing.T) {}
