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
