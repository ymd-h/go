// slices package
//
// Generic Class base implementation for golang.org/x/exp/slices
package slices

import (
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

type (
	ordered = constraints.Ordered

	ISlice[E any] interface {
		Get(int) E
		Append(E)
		Size() int
		BynarySearchFunc(E, func(E, E) int) (int, bool)
		Clip()
		Clone() ISlice[E]
		CompactFunc(func (E, E) bool)
		CompareFunc(ISlice[E], func (E, E) int) int
		ContainsFunc(func (E) bool) bool
		Delete(int, int)
		EqualFunc(ISlice[E], func (E, E) bool) bool
		Grow(int)
		IndexFunc(func (E) bool) int
		Insert(int, ...E)
		IsSortedFunc(func (E, E) bool) bool
		Replace(int, int, ...E)
		SortFunc(func (E, E) bool)
		SortStableFunc(func (E, E) bool)
		TryComparable() (IComparableSlice[E], bool)
		TryOrdered() (IOrderedSlice[E], bool)
	}
	IComparableSlice[E any] interface {
		ISlice[E]
		Compact()
		Contains(E) bool
		Equal(IComparableSlice[E]) bool
		Index(E) int
	}
	IOrderedSlice[E any] interface {
		IComparableSlice[E]
		BynarySearch(E) (int, bool)
		Compare(IOrderedSlice[E]) int
		IsSorted() bool
		Sort()
	}

	Slice[E any] struct {
		item []E
	}
	ComparableSlice[E comparable] struct {
		Slice[E]
	}
	OrderedSlice[E ordered] struct {
		ComparableSlice[E]
	}
)


func NewSlice[E any]() ISlice[E] {
	var e E
	switch any(e).(type) {
	case bool:
		return NewComparableSlice[bool]().(ISlice[E])
	case string:
		return NewComparableSlice[string]().(ISlice[E])
	case int:
		return NewComparableSlice[int]().(ISlice[E])
	case int8:
		return NewComparableSlice[int8]().(ISlice[E])
	case int16:
		return NewComparableSlice[int16]().(ISlice[E])
	case int32:
		return NewComparableSlice[int32]().(ISlice[E])
	case int64:
		return NewComparableSlice[int64]().(ISlice[E])
	case uint:
		return NewComparableSlice[uint]().(ISlice[E])
	case uint8:
		return NewComparableSlice[uint8]().(ISlice[E])
	case uint16:
		return NewComparableSlice[uint16]().(ISlice[E])
	case uint32:
		return NewComparableSlice[uint32]().(ISlice[E])
	case uint64:
		return NewComparableSlice[uint64]().(ISlice[E])
	case uintptr:
		return NewComparableSlice[uintptr]().(ISlice[E])
	case float32:
		return NewComparableSlice[float32]().(ISlice[E])
	case float64:
		return NewComparableSlice[float64]().(ISlice[E])
	case complex64:
		return NewComparableSlice[complex64]().(ISlice[E])
	case complex128:
		return NewComparableSlice[complex128]().(ISlice[E])
	default:
		return &Slice[E]{
			item: make([]E, 0),
		}
	}

}
func NewComparableSlice[E comparable]() IComparableSlice[E] {
	var e E
	switch any(e).(type) {
	case string:
		return NewOrderedSlice[string]().(IComparableSlice[E])
	case int:
		return NewOrderedSlice[int]().(IComparableSlice[E])
	case int8:
		return NewOrderedSlice[int8]().(IComparableSlice[E])
	case int16:
		return NewOrderedSlice[int16]().(IComparableSlice[E])
	case int32:
		return NewOrderedSlice[int32]().(IComparableSlice[E])
	case int64:
		return NewOrderedSlice[int64]().(IComparableSlice[E])
	case uint:
		return NewOrderedSlice[uint]().(IComparableSlice[E])
	case uint8:
		return NewOrderedSlice[uint8]().(IComparableSlice[E])
	case uint16:
		return NewOrderedSlice[uint16]().(IComparableSlice[E])
	case uint32:
		return NewOrderedSlice[uint32]().(IComparableSlice[E])
	case uint64:
		return NewOrderedSlice[uint64]().(IComparableSlice[E])
	case float32:
		return NewOrderedSlice[float32]().(IComparableSlice[E])
	case float64:
		return NewOrderedSlice[float64]().(IComparableSlice[E])
	default:
		return &ComparableSlice[E]{
			Slice[E]{
				item: make([]E, 0),
			},
		}
	}

}
func NewOrderedSlice[E ordered]() IOrderedSlice[E] {
	return &OrderedSlice[E]{
		ComparableSlice[E]{
			Slice[E]{
				item: make([]E, 0),
			},
		},
	}
}


func (s *Slice[E]) Get(i int) E {
	return s.item[i]
}

func (s *Slice[E]) Append(e E) {
	s.item = append(s.item, e)
}

func (s *Slice[E]) Size() int {
	return len(s.item)
}

func search(n int, gt func(int) bool) int {
	start, end := 0, n
	for start < end {
		check := start + (end - start) / 2
		if gt(check) {
			end = check
		} else {
			start = check + 1
		}
	}
	return start
}

func (s *Slice[E]) BynarySearchFunc(target E, cmp func(E, E) int) (int, bool) {
	p := search(len(s.item), func(i int) bool { return cmp(s.item[i], target) >= 0 })
	if (p >= len(s.item)) || (cmp(s.item[p], target) != 0) {
		return p, false
	} else {
		return p, true
	}
}

func (s *Slice[E]) Clip() {
	n := len(s.item)
	s.item = s.item[:n:n]
}

func (s *Slice[E]) Clone() ISlice[E] {
	return &Slice[E]{
		item: append(make([]E, 0, len(s.item)), s.item...),
	}
}

func (s *Slice[E]) CompactFunc(eq func(E, E) bool) {
	if len(s.item) < 2 {
		return
	}

	to := 1
	last := s.item[0]

	for _, v := range s.item[1:] {
		if !eq(v, last) {
			s.item[to] = v
			to += 1
			last = v
		}
	}

	s.item = s.item[:to]
}

func (s *Slice[E]) CompareFunc(o ISlice[E], cmp func(E, E) int) int {
	return CompareFunc[E, E](s, o, cmp)
}

func (s *Slice[E]) ContainsFunc(cond func(E) bool) bool {
	for _, v := range s.item {
		if cond(v) {
			return true
		}
	}
	return false
}


func (s *Slice[E]) Delete(delStart, delEnd int) {
	n := len(s.item)
	if delStart >= n {
		return
	}
	if delEnd >= n {
		s.item = s.item[:delStart]
		return
	}

	s.item = append(s.item[:delStart], s.item[delEnd:]...)
}


func (s *Slice[E]) EqualFunc(o ISlice[E], eq func (E, E) bool) bool {
	if len(s.item) != o.Size() {
		return false
	}
	for i, v := range s.item {
		if !eq(v, o.Get(i)) {
			return false
		}
	}
	return true
}


func (s *Slice[E]) Grow(n int) {
	if n <= 0 {
		return
	}
	want := len(s.item) + n
	if want > cap(s.item) {
		s.item = append(make([]E, 0, want), s.item...)
	}
}


func (s *Slice[E]) IndexFunc(cond func (E) bool) int {
	for i, v := range s.item {
		if cond(v) {
			return i
		}
	}
	return -1
}


func (s *Slice[E]) Insert(i int, elems ...E) {
	n := len(elems)
	need := len(s.item) + n
	if cap(s.item) > need {
		copy(s.item[i+n:need], s.item[i:])
		copy(s.item[i:], elems)
		return
	}

	item := make([]E, need)
	copy(item, s.item[:i])
	copy(item[i:], elems)
	copy(item[i+n:], s.item[i:])
	s.item = item
}


func (s *Slice[E]) IsSortedFunc(less func (E, E) bool) bool {
	n := len(s.item)
	if n < 2 {
		return true
	}

	for i := 1; i < n; i++ {
		if less(s.item[i], s.item[i-1]) {
			return false
		}
	}

	return true
}


func (s *Slice[E]) Replace(start, end int, elems ...E) {
	slen := len(s.item)
	if (start >= end) || (start >= slen) {
		return
	}
	if end > slen {
		end = slen
	}

	n := end - start
	switch {
	case n > len(elems):
		end = start + len(elems)
	case n < len(elems):
		elems = elems[:n]
	}

	copy(s.item[start:end], elems)
}

func (s *Slice[E]) SortFunc(less func(E, E) bool) {
	slices.SortFunc(s.item, less)
}


func (s *Slice[E]) SortStableFunc(less func(E, E) bool) {
	slices.SortStableFunc(s.item, less)
}


func (s *Slice[E]) TryComparable() (IComparableSlice[E], bool) {
	return nil, false
}

func (s *Slice[E]) TryOrdered() (IOrderedSlice[E], bool) {
	return nil, false
}


func (s *ComparableSlice[E]) Clone() ISlice[E] {
	return &ComparableSlice[E]{
		Slice[E]{
			item: append(make([]E, 0, len(s.item)), s.item...),
		},
	}
}

func (s *ComparableSlice[E]) Compact() {
	if len(s.item) < 2 {
		return
	}

	to := 1
	last := s.item[0]

	for _, v := range s.item[1:] {
		if v != last {
			s.item[to] = v
			to += 1
			last = v
		}
	}

	s.item = s.item[:to]
}

func (s *ComparableSlice[E]) Contains(e E) bool {
	for _, v := range s.item {
		if v == e {
			return true
		}
	}
	return false
}


func (s *ComparableSlice[E]) Equal(o IComparableSlice[E]) bool {
	if len(s.item) != o.Size() {
		return false
	}
	for i, v := range s.item {
		if v != o.Get(i) {
			return false
		}
	}
	return true
}


func (s *ComparableSlice[E]) Index(e E) int {
	for i, v := range s.item {
		if v == e {
			return i
		}
	}
	return -1
}


func (s *ComparableSlice[E]) TryComparable() (IComparableSlice[E], bool) {
	return any(s).(IComparableSlice[E]), true
}

func (s *ComparableSlice[E]) TryOrdered() (IOrderedSlice[E], bool) {
	return nil, false
}



func (s *OrderedSlice[E]) BynarySearch(target E) (int, bool) {
	p := search(len(s.item), func(i int) bool { return s.item[i] >= target })
	if (p >= len(s.item)) || (s.item[p] != target) {
		return p, false
	} else {
		return p, true
	}
}


func (s *OrderedSlice[E]) Clone() ISlice[E] {
	return &OrderedSlice[E]{
		ComparableSlice[E]{
			Slice[E]{
				item: append(make([]E, 0, len(s.item)), s.item...),
			},
		},
	}
}


func (s *OrderedSlice[E]) Compare(o IOrderedSlice[E]) int {
	olen := o.Size()
	for i, v1 := range s.item {
		if i >= olen {
			return +1
		}

		v2 := o.Get(i)
		switch {
		case v1 < v2:
			return -1
		case v1 > v2:
			return +1
		}
	}
	if s.Size() < olen {
		return -1
	}
	return 0
}


func (s *OrderedSlice[E]) IsSorted() bool {
	n := len(s.item)
	if n < 2 {
		return true
	}

	for i := 1; i < n; i++ {
		if s.item[i] < s.item[i-1] {
			return false
		}
	}

	return true
}


func (s *OrderedSlice[E]) Sort() {
	slices.Sort(s.item)
}


func (s *OrderedSlice[E]) TryComparable() (IComparableSlice[E], bool) {
	return any(s).(IComparableSlice[E]), true
}

func (s *OrderedSlice[E]) TryOrdered() (IOrderedSlice[E], bool) {
	return any(s).(IOrderedSlice[E]), true
}



func CompareFunc[E1, E2 any](s1 ISlice[E1], s2 ISlice[E2], cmp func(E1, E2) int) int {
	s1len := s1.Size()
	s2len := s2.Size()

	for i := 0; i < s1len; i++ {
		if i >= s2len {
			return +1
		}

		v1 := s1.Get(i)
		v2 := s2.Get(i)
		if c := cmp(v1, v2); c != 0 {
			return c
		}
	}
	if s1len < s2len {
		return -1
	}
	return 0
}


func EqualFunc[E1, E2 any](s1 ISlice[E1], s2 ISlice[E2], eq func(E1, E2) bool) bool {
	s1len := s1.Size()
	if s1len !=s2.Size() {
		return false
	}

	for i := 0; i < s1len; i++ {
		if !eq(s1.Get(i), s2.Get(i)) {
			return false
		}
	}

	return true
}
