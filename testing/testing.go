package testing

import (
	"testing"
)

type (
	testData[T any] struct {
		name string
		data T
	}
	Test[T any] struct {
		t *testing.T
		tests []testData[T]
	}
)

func NewTest[T any](t *testing.T) *Test[T] {
	return &Test[T]{
		t: t,
		tests: make([]testData[T], 0),
	}
}

func (p *Test[T]) Add(name string, data T) *Test[T] {
	p.tests = append(p.tests, testData[T]{ name: name, data: data })
	return p
}

func (p *Test[T]) Run(f func(*testing.T, T)) *Test[T] {
	for _, tt := range p.tests {
		p.t.Run(tt.name, func(t *testing.T) {
			f(t, tt.data)
		})
	}
	return p
}




func AssertEqual[T comparable](t *testing.T, got T, want T) {
	if got != want {
		t.Errorf("%v != %v\n", got, want)
	}
}
