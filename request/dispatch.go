package request

import (
	"fmt"
)

type (
	// Interface for Dispatch Item
	IDispatchItem interface {
		StatusCode() int
		Response() any
	}

	// Dispatch Item
	DispatchItem[T any] struct {
		Code int
	}

	// Response Dispatcher class
	ResponseDispatcher struct {
		rule map[int] IDispatchItem
	}
)


func (i DispatchItem[T]) StatusCode() {
	return i.Code
}

func (i DispatchItem[T]) Response() any {
	var v T
	return &v
}

func NewResponseDispatcher(items ...IDispatchItem) *ResponseDispatcher {
	d := &ResponseDispatcher{rule: make(map[int] IDispatchItem, len(items))}

	d.Add(items...)

	return d
}

func (d *ResponseDispatcher) Add(items ...IDispatchItem) {
	for _, i := range items {
		d.rule[i.StatusCode()] = i
	}
}

func (d *ResponseDispatcher) Delete(codes ...int) {
	for _, i := range codes {
		delete(d.rule, i)
	}
}

func (d *ResponseDispatcher) Dispatch(code int) (any, error) {
	if r, ok := d.rule[code]; ok {
		return r.Response(), nil
	}

	return nil, fmt.Errorf("Fail to Dispatch StatusCode: %d", code)
}
