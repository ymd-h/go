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
		// Get value at specified index
		//
		// # Arguments
		// * `i`: `int` - Index
		//
		// # Returns
		// * `E` - Internal value
		Get(i int) E

		// Set value at specified index
		//
		// # Arguments
		// * `i`: `int` - Index
		// * `e`: `E` - Value to be set
		Set(i int, e E)

		// Append element(s)
		//
		// # Arguments
		// `elems`: `...E` - Elements to be added
		Append(elems ...E)
		Size() int
		BinarySearchFunc(E, func(E, E) int) (int, bool)
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
		BinarySearch(E) (int, bool)
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
	return NewSliceFrom(make([]E, 0))
}

func NewComparableSlice[E comparable]() IComparableSlice[E] {
	return NewComparableSliceFrom(make([]E, 0))
}

func NewOrderedSlice[E ordered]() IOrderedSlice[E] {
	return NewOrderedSliceFrom(make([]E, 0))
}


func NewSliceFrom[E any, S []E](s S) ISlice[E] {
	switch t := any(s).(type) {
	case []bool:
		return NewComparableSliceFrom[bool](t).(ISlice[E])
	case []string:
		return NewComparableSliceFrom[string](t).(ISlice[E])
	case []int:
		return NewComparableSliceFrom[int](t).(ISlice[E])
	case []int8:
		return NewComparableSliceFrom[int8](t).(ISlice[E])
	case []int16:
		return NewComparableSliceFrom[int16](t).(ISlice[E])
	case []int32:
		return NewComparableSliceFrom[int32](t).(ISlice[E])
	case []int64:
		return NewComparableSliceFrom[int64](t).(ISlice[E])
	case []uint:
		return NewComparableSliceFrom[uint](t).(ISlice[E])
	case []uint8:
		return NewComparableSliceFrom[uint8](t).(ISlice[E])
	case []uint16:
		return NewComparableSliceFrom[uint16](t).(ISlice[E])
	case []uint32:
		return NewComparableSliceFrom[uint32](t).(ISlice[E])
	case []uint64:
		return NewComparableSliceFrom[uint64](t).(ISlice[E])
	case []uintptr:
		return NewComparableSliceFrom[uintptr](t).(ISlice[E])
	case []float32:
		return NewComparableSliceFrom[float32](t).(ISlice[E])
	case []float64:
		return NewComparableSliceFrom[float64](t).(ISlice[E])
	case []complex64:
		return NewComparableSliceFrom[complex64](t).(ISlice[E])
	case []complex128:
		return NewComparableSliceFrom[complex128](t).(ISlice[E])
	default:
		return &Slice[E]{
			item: s,
		}
	}
}
func NewComparableSliceFrom[E comparable, S []E](s S) IComparableSlice[E] {
	switch t := any(s).(type) {
	case []string:
		return NewOrderedSliceFrom[string](t).(IComparableSlice[E])
	case []int:
		return NewOrderedSliceFrom[int](t).(IComparableSlice[E])
	case []int8:
		return NewOrderedSliceFrom[int8](t).(IComparableSlice[E])
	case []int16:
		return NewOrderedSliceFrom[int16](t).(IComparableSlice[E])
	case []int32:
		return NewOrderedSliceFrom[int32](t).(IComparableSlice[E])
	case []int64:
		return NewOrderedSliceFrom[int64](t).(IComparableSlice[E])
	case []uint:
		return NewOrderedSliceFrom[uint](t).(IComparableSlice[E])
	case []uint8:
		return NewOrderedSliceFrom[uint8](t).(IComparableSlice[E])
	case []uint16:
		return NewOrderedSliceFrom[uint16](t).(IComparableSlice[E])
	case []uint32:
		return NewOrderedSliceFrom[uint32](t).(IComparableSlice[E])
	case []uint64:
		return NewOrderedSliceFrom[uint64](t).(IComparableSlice[E])
	case []float32:
		return NewOrderedSliceFrom[float32](t).(IComparableSlice[E])
	case []float64:
		return NewOrderedSliceFrom[float64](t).(IComparableSlice[E])
	default:
		return &ComparableSlice[E]{
			Slice[E]{
				item: s,
			},
		}
	}
}
func NewOrderedSliceFrom[E ordered, S []E](s S) IOrderedSlice[E] {
	return &OrderedSlice[E]{
		ComparableSlice[E]{
			Slice[E]{
				item: s,
			},
		},
	}
}

func (s *Slice[E]) Get(i int) E {
	return s.item[i]
}

func (s *Slice[E]) Set(i int, e E) {
	s.item[i] = e
}

func (s *Slice[E]) Append(elems ...E) {
	s.item = append(s.item, elems...)
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

func (s *Slice[E]) BinarySearchFunc(target E, cmp func(E, E) int) (int, bool) {
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



func (s *OrderedSlice[E]) BinarySearch(target E) (int, bool) {
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
