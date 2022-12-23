package sets

type (
	ISet[T any] interface {
		Add(T)
		Remove(T)
		Has(T) bool
		Size() int
		Copy(ISet[T])
		Clone() ISet[T]
		Equal(ISet[T]) bool
		Difference(ISet[T]) ISet[T]
		Intersection(ISet[T]) ISet[T]
		Union(ISet[T]) ISet[T]
		SubsetOf(ISet[T]) bool
		ProperSubsetOf(ISet[T]) bool
	}
	Set[T comparable] struct {
		item map[T]struct{}
	}

	// Compatible with "github.com/ymd-h/go/slices"
	ISlice[T any] interface {
		Get(int) T
		Size() int
	}
)

func New[T comparable]() ISet[T] {
	return &Set[T]{ item: make(map[T]struct{}, 0) }
}

func FromSlice[T comparable](s []T) ISet[T] {
	set := New[T]()
	for _, v := range s {
		set.Add(v)
	}
	return set
}

func FromISlice[T comparable](s ISlice[T]) ISet[T] {
	set := New[T]()
	n := s.Size()
	for i := 0; i < n; i++ {
		set.Add(s.Get(i))
	}
	return set
}

func (p *Set[T]) Add(v T) {
	p.item[v] = struct{}{}
}

func (p *Set[T]) Remove(v T) {
	delete(p.item, v)
}

func (p *Set[T]) Has(v T) bool {
	_, ok := p.item[v]
	return ok
}

func (p *Set[T]) Size() int {
	return len(p.item)
}

func (p *Set[T]) Copy(s ISet[T]) {
	for v, _ := range p.item {
		s.Add(v)
	}
}

func (p *Set[T]) Clone() ISet[T] {
	s := New[T]()
	p.Copy(s)
	return s
}

func (p *Set[T]) Equal(s ISet[T]) bool {
	if len(p.item) != s.Size() {
		return false
	}
	for v, _ := range p.item {
		if !s.Has(v) {
			return false
		}
	}
	return true
}

func (p *Set[T]) Difference(s ISet[T]) ISet[T] {
	d := s.Clone()
	for v, _ := range p.item {
		if s.Has(v) {
			d.Remove(v)
		} else {
			d.Add(v)
		}
	}
	return d
}

func (p *Set[T]) Intersection(s ISet[T]) ISet[T] {
	i := New[T]()
	for v, _ := range p.item {
		if s.Has(v) {
			i.Add(v)
		}
	}
	return i
}

func (p *Set[T]) Union(s ISet[T]) ISet[T] {
	u := p.Clone()
	s.Copy(u)
	return u
}

func (p *Set[T]) SubsetOf(s ISet[T]) bool {
	if len(p.item) > s.Size() {
		return false
	}
	for v, _ := range p.item {
		if !s.Has(v) {
			return false
		}
	}
	return true
}

func (p *Set[T]) ProperSubsetOf(s ISet[T]) bool {
	if len(p.item) >= s.Size() {
		return false
	}
	for v, _ := range p.item {
		if !s.Has(v) {
			return false
		}
	}
	return true
}
