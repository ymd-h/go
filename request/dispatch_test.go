package request

import (
	"testing"
)


func TestDispatcher(t *testing.T){
	type (
		A struct {}
		B struct {}
		C struct {}
	)

	d := NewResponseDispatcher(
		DispatchItem[A]{200},
		DispatchItem[B]{300},
		DispatchItem[C]{400},
	)

	a, err := d.Dispatch(200)
	if err != nil {
		t.Errorf("Fail: %v\n", err)
		return
	}
	if _, ok := a.(*A); !ok {
		t.Errorf("Fail: %T\n", a)
		return
	}

	b, err := d.Dispatch(300)
	if err != nil {
		t.Errorf("Fail: %v\n", err)
		return
	}
	if _, ok := b.(*B); !ok {
		t.Errorf("Fail: %T\n", b)
		return
	}

	c, err := d.Dispatch(400)
	if err != nil {
		t.Errorf("Fail: %v\n", err)
		return
	}
	if _, ok := c.(*C); !ok {
		t.Errorf("Fail: %T\n", c)
		return
	}

	_, err = d.Dispatch(500)
	if err == nil {
		t.Errorf("Must fail\n")
		return
	}
}
