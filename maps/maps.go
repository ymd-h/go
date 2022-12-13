// maps package
//
// Generics Class base implementation for golang.org/x/exp/maps
package maps

type (
	IMap[K comparable, V any] interface {
		Get(K) (V, bool)
		Set(K, V)
		Keys() []K
		Values() []V
		Size() int
		Clear()
		Clone() IMap[K, V]
		Copy(IMap[K, V]) IMap[K, V]
		Delete(K)
		DeleteFunc(func (K, V) bool)
		EqualFunc(IMap[K, V], func(V, V) bool) bool
		TryComparable() (IComparableMap[K, V], bool)
	}
	IComparableMap[K comparable, V any] interface {
		IMap[K, V]
		Equal(IComparableMap[K, V]) bool
	}

	Map[K comparable, V any] struct {
		item map[K]V
	}
	ComparableMap[K, V comparable] struct {
		Map[K, V]
	}
)

func NewMap[K comparable, V any]() IMap[K, V] {
	var v V
	switch any(v).(type) {
	case bool:
		return NewComparableMap[K, bool]().(IMap[K, V])
	case string:
		return NewComparableMap[K, string]().(IMap[K, V])
	case int:
		return NewComparableMap[K, int]().(IMap[K, V])
	case int8:
		return NewComparableMap[K, int8]().(IMap[K, V])
	case int16:
		return NewComparableMap[K, int16]().(IMap[K, V])
	case int32:
		return NewComparableMap[K, int32]().(IMap[K, V])
	case int64:
		return NewComparableMap[K, int64]().(IMap[K, V])
	case uint:
		return NewComparableMap[K, uint]().(IMap[K, V])
	case uint8:
		return NewComparableMap[K, uint8]().(IMap[K, V])
	case uint16:
		return NewComparableMap[K, uint16]().(IMap[K, V])
	case uint32:
		return NewComparableMap[K, uint32]().(IMap[K, V])
	case uint64:
		return NewComparableMap[K, uint64]().(IMap[K, V])
	case uintptr:
		return NewComparableMap[K, uintptr]().(IMap[K, V])
	case float32:
		return NewComparableMap[K, float32]().(IMap[K, V])
	case float64:
		return NewComparableMap[K, float64]().(IMap[K, V])
	case complex64:
		return NewComparableMap[K, complex64]().(IMap[K, V])
	case complex128:
		return NewComparableMap[K, complex128]().(IMap[K, V])
	default:
		return &Map[K, V]{
			item: make(map[K]V, 0),
		}
	}
}

func NewMapFrom[K comparable, V any](m map[K]V) IMap[K, V] {
	switch t := any(m).(type) {
	case map[K]bool:
		return NewComparableMapFrom(t).(IMap[K, V])
	case map[K]string:
		return NewComparableMapFrom(t).(IMap[K, V])
	case map[K]int:
		return NewComparableMapFrom(t).(IMap[K, V])
	case map[K]int8:
		return NewComparableMapFrom(t).(IMap[K, V])
	case map[K]int16:
		return NewComparableMapFrom(t).(IMap[K, V])
	case map[K]int32:
		return NewComparableMapFrom(t).(IMap[K, V])
	case map[K]int64:
		return NewComparableMapFrom(t).(IMap[K, V])
	case map[K]uint:
		return NewComparableMapFrom(t).(IMap[K, V])
	case map[K]uint8:
		return NewComparableMapFrom(t).(IMap[K, V])
	case map[K]uint16:
		return NewComparableMapFrom(t).(IMap[K, V])
	case map[K]uint32:
		return NewComparableMapFrom(t).(IMap[K, V])
	case map[K]uint64:
		return NewComparableMapFrom(t).(IMap[K, V])
	case map[K]uintptr:
		return NewComparableMapFrom(t).(IMap[K, V])
	case map[K]float32:
		return NewComparableMapFrom(t).(IMap[K, V])
	case map[K]float64:
		return NewComparableMapFrom(t).(IMap[K, V])
	case map[K]complex64:
		return NewComparableMapFrom(t).(IMap[K, V])
	case map[K]complex128:
		return NewComparableMapFrom(t).(IMap[K, V])
	default:
		return &Map[K, V]{
			item: m,
		}
	}
}

func NewComparableMap[K, V comparable]() IComparableMap[K, V] {
	return &ComparableMap[K, V]{
		Map[K,V]{
			item: make(map[K]V, 0),
		},
	}
}

func NewComparableMapFrom[K, V comparable](m map[K]V) IComparableMap[K, V] {
	return &ComparableMap[K, V]{
		Map[K, V]{
			item: m,
		},
	}
}

func (m *Map[K, V]) Get(k K) (V, bool) {
	v, ok := m.item[k]
	return v, ok
}

func (m *Map[K, V]) Set(k K, v V) {
	m.item[k] = v
}

func (m *Map[K, V]) Keys() []K {
	keys := make([]K, len(m.item))
	for k, _ := range m.item {
		keys = append(keys, k)
	}
	return keys
}

func (m *Map[K, V]) Values() []V {
	values := make([]V, len(m.item))
	for _, v := range m.item {
		values = append(values, v)
	}
	return values
}

func (m *Map[K, V]) Size() int {
	return len(m.item)
}

func (m *Map[K, V]) Clear() {
	for k := range m.item {
		delete(m.item, k)
	}
}

func (m *Map[K, V]) Clone() IMap[K, V] {
	item := make(map[K]V, len(m.item))
	for k, v := range m.item {
		item[k] = v
	}
	return &Map[K, V]{ item: item }
}

func (m *Map[K, V]) Copy(o IMap[K, V]) IMap[K, V] {
	for k, v := range m.item {
		o.Set(k, v)
	}
	return o
}

func (m *Map[K, V]) Delete(k K) {
	delete(m.item, k)
}

func (m *Map[K, V]) DeleteFunc(del func(K, V) bool) {
	for k, v := range m.item {
		if del(k, v) {
			delete(m.item, k)
		}
	}
}

func (m *Map[K, V]) EqualFunc(o IMap[K, V], eq func(V, V) bool) bool {
	return EqualFunc[K, V, V](m, o, eq)
}

func (m *Map[K, V]) TryComparable() (IComparableMap[K, V], bool) {
	return nil, false
}

func (m *ComparableMap[K, V]) Clone() IMap[K, V] {
	item := make(map[K]V, len(m.item))
	for k, v := range m.item {
		item[k] = v
	}
	return &ComparableMap[K, V]{
		Map[K, V]{
			item: item,
		},
	}
}

func (m *ComparableMap[K, V]) TryComparable() (IComparableMap[K, V], bool) {
	return m, true
}

func (m *ComparableMap[K, V]) Equal(o IComparableMap[K, V]) bool {
	if len(m.item) != o.Size() {
		return false
	}

	for k, v := range m.item {
		if ov, ok := o.Get(k); (!ok) || (v != ov) {
			return false
		}
	}

	return true
}

func Equal[K, V comparable](m1, m2 IComparableMap[K, V]) bool {
	return m1.Equal(m2)
}


func EqualFunc[K comparable, V1, V2 any](m1 IMap[K, V1], m2 IMap[K, V2], eq func (V1, V2) bool) bool {
	if m1.Size() != m2.Size() {
		return false
	}

	var key K
	for k := range m1.Keys() {
		// # Issue
		// Go 1.19 compiler cannot estimate type correctly, and fail to build;
		// `cannot use k (variable of type int) as type K in argument to m1.Get` etc.
		// As a workaround, we dynamically assert type
		key = any(k).(K)
		v1, _ := m1.Get(key)
		if v2, ok := m2.Get(key); (!ok) || (!eq(v1, v2)) {
			return false
		}
	}

	return true
}
