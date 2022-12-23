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
