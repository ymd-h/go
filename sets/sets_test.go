package sets

import (
	"testing"

	"github.com/ymd-h/go/slices"
	y "github.com/ymd-h/go/testing"
)

func TestNewSet(t *testing.T) {
	type test struct {
		args ISet[int]
	}
	y.NewTest[test](t).
		Add("New", test{
			args: New[int](),
		}).
		Add("FromSlice", test{
			args: FromSlice[int]([]int{}),
		}).
		Add("FromISlice", test{
			args: FromISlice[int](slices.NewSlice[int]()),
		}).
		Run(func(_ *testing.T, data test) {
			switch data.args.(type) {
			case *Set[int]:
			default:
				t.Errorf("Type mismatch: %T != *Set[int]\n", data.args)
			}
		})
}


func TestAdd(t *testing.T) {
	type test struct {
		set ISet[int]
		key int
		wantSize int
	}
	y.NewTest[test](t).
		Add("simple", test{
			set: New[int](),
			key: 1,
			wantSize: 1,
		}).
		Add("duplicated", test{
			set: FromSlice[int]([]int{1}),
			key: 1,
			wantSize: 1,
		}).
		Run(func(_ *testing.T, data test) {
			data.set.Add(data.key)
			y.AssertEqual(t, data.set.Size(), data.wantSize)
			y.AssertEqual(t, data.set.Has(data.key), true)
		})
}


func TestRemove(t *testing.T) {
	type test struct {
		set ISet[string]
		key string
		wantSize int
	}
	y.NewTest[test](t).
		Add("simple", test{
			set: FromSlice[string]([]string{"abc"}),
			key: "abc",
			wantSize: 0,
		}).
		Add("empty", test{
			set: New[string](),
			key: "abc",
			wantSize: 0,
		}).
		Run(func(_ *testing.T, data test) {
			data.set.Remove(data.key)
			y.AssertEqual(t, data.set.Size(), data.wantSize)
		})
}


func TestHas(t *testing.T) {
	type test struct {
		set ISet[float32]
		key float32
		want bool
	}
	y.NewTest[test](t).
		Add("has", test{
			set: FromSlice([]float32{0.5}),
			key: 0.5,
			want: true,
		}).
		Add("hasn't", test{
			set: New[float32](),
			key: 0.5,
			want: false,
		}).
		Run(func(_ *testing.T, data test) {
			y.AssertEqual(t, data.set.Has(data.key), data.want)
		})
}


func TestSize(t *testing.T) {
	type test struct {
		set ISet[int]
		want int
	}
	y.NewTest[test](t).
		Add("empty", test{
			set: New[int](),
			want: 0,
		}).
		Add("some", test{
			set: FromSlice([]int{1, 2}),
			want: 2,
		}).
		Run(func(_ *testing.T, data test) {
			y.AssertEqual(t, data.set.Size(), data.want)
		})
}

func TestToSlice(t *testing.T) {
	type test struct{
		set ISet[int]
		want []int
	}
	y.NewTest[test](t).
		Add("empty", test{
			set: New[int](),
			want: []int{},
		}).
		Add("some", test{
			set: FromSlice([]int{1, 2}),
			want: []int{1, 2},
		}).
		Run(func(_ *testing.T, data test) {
			s := data.set.ToSlice()
			y.AssertEqual(t, len(s), len(data.want))
			for _, ss := range s {
				y.AssertIsIn(t, ss, data.want)
			}
		})
}


func TestCopy(t *testing.T) {
	type test struct {
		setFrom ISet[int]
		setTo ISet[int]
		wantSize int
	}
	y.NewTest[test](t).
		Add("no overlap", test{
			setFrom: FromSlice([]int{1, 2, 3}),
			setTo: New[int](),
			wantSize: 3,
		}).
		Add("overlap", test{
			setFrom: FromSlice([]int{1, 2, 3}),
			setTo: FromSlice([]int{3, 4, 5}),
			wantSize: 5,
		}).
		Run(func(_ *testing.T, data test) {
			data.setFrom.Copy(data.setTo)
			y.AssertEqual(t, data.setTo.Size(), data.wantSize)
			for _, v := range data.setFrom.ToSlice() {
				y.AssertEqual(t, data.setTo.Has(v), true)
			}
		})
}


func TestClone(t *testing.T) {
	type test struct {
		set ISet[int]
	}
	y.NewTest[test](t).
		Add("simple", test{
			set: FromSlice([]int{1, 2, 3}),
		}).
		Run(func(_ *testing.T, data test) {
			c := data.set.Clone()
			y.AssertEqual(t, c.Size(), data.set.Size())
			for _, e := range c.ToSlice() {
				y.AssertEqual(t, data.set.Has(e), true)
			}
			data.set.Add(999999)
			y.AssertEqual(t, c.Size() + 1, data.set.Size())
		})
}


func TestEqual(t *testing.T) {
	type test struct {
		setA ISet[int]
		setB ISet[int]
		want bool
	}
	y.NewTest[test](t).
		Add("true", test{
			setA: FromSlice([]int{1, 2, 3}),
			setB: FromSlice([]int{2, 3, 1}),
			want: true,
		}).
		Add("false", test{
			setA: FromSlice([]int{1, 2}),
			setB: FromSlice([]int{2, 3}),
			want: false,
		}).
		Run(func(_ *testing.T, data test) {
			y.AssertEqual(t, data.setA.Equal(data.setB), data.want)
		})
}


func TestDifference(t *testing.T) {
	type test struct {
		setA ISet[int]
		setB ISet[int]
		want ISet[int]
	}
	y.NewTest[test](t).
		Add("simple", test{
			setA: FromSlice([]int{1, 2}),
			setB: FromSlice([]int{2, 3}),
			want: FromSlice([]int{1, 3}),
		}).
		Add("same", test{
			setA: FromSlice([]int{1, 2, 3}),
			setB: FromSlice([]int{1, 2, 3}),
			want: New[int](),
		}).
		Run(func(_ *testing.T, data test) {
			d := data.setA.Difference(data.setB)
			y.AssertEqual(t, d.Equal(data.want), true)
		})
}
