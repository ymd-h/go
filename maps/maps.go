// maps package
//
// Generics Class base implementation for golang.org/x/exp/maps
package maps

type (
	// Interface for Generic Map
	IMap[K comparable, V any] interface {
		// Get value at key if it exists
		//
		// # Arguments
		// * `key`: `K` - Key
		//
		// # Returns
		// * `V` - Value
		// * `bool` - Success
		Get(key K) (V, bool)

		// Set value at key
		//
		// # Arguments
		// * `key`: `K` - Key
		// * `value`: `V` - Value
		Set(key K, value V)

		// Get keys of map
		//
		// # Returns
		// * `[]K` - Slice of keys
		Keys() []K

		// Get Values of map
		//
		// # Returns
		// * `[]V` - Slice of values
		Values() []V

		// Get size of map
		//
		// # Returns
		// * `int` - Size
		Size() int

		// Clear map
		Clear()

		// Clone map
		//
		// # Returns
		// * `IMap[K, V]` - Cloned map
		//
		// # Warning
		// This is shallow copy
		Clone() IMap[K, V]

		// Copy to another map
		//
		// # Arguments
		// `other`: `IMap[K, V]` - Copy destination
		//
		// # Returns
		// * `IMap[K, V]` - `other`
		Copy(other IMap[K, V]) IMap[K, V]

		// Delete value at a key
		//
		// # Arguments
		// * `key`: `K` - Key to be deleted
		Delete(key K)

		// Delete values which meets condition
		//
		// # Arguments
		// * `del`: `func(K, V) bool` - Delete condition
		DeleteFunc(del func (K, V) bool)

		// Equal check with function
		//
		// # Arguments
		// * `other`: `IMap[K, V]` - Comparison target
		//
		// # Returns
		// * `eq`: `func(V, V) bool` - Equality function
		//
		// # Notes
		// Since method cannot be generics, `other` is limitied to
		// the same type paramter `[K, V]`.
		// Free function version can accept different value type, too.
		EqualFunc(other IMap[K, V], eq func(V, V) bool) bool

		// Type assertion to IComparable interface
		//
		// # Returns
		// * `IComparable[K, V]` - Comparable Map if possible
		// * `bool` - Success
		TryComparable() (IComparableMap[K, V], bool)
	}

	// Interface for Comparable Map
	//
	// # Notes
	// Interface is defined over `[K comparable, V any]`
	// to support `IMap[K, V].TryComparable()`,
	// however, implementation supports only `[K, V comparable]`
	IComparableMap[K comparable, V any] interface {
		IMap[K, V]

		// Equality check
		//
		// # Arguments
		// * `other`: `IComparableMap[K, V]` - Map to be checked
		//
		// # Returns
		// * `bool` - Whether they are equal or not
		Equal(other IComparableMap[K, V]) bool
	}

	// Class implements `IMap[K, V]`
	Map[K comparable, V any] struct {
		item map[K]V
	}

	// Class implements `IComparableMap[K, V]`
	ComparableMap[K, V comparable] struct {
		Map[K, V]
	}
)

// Create new map
//
// # Returns
// * `IMap[K, V]` - Created map
//
// # Notes
// For some predefined primitive types, we detect and create `ComparableMap[K, V]`.
// This mechanism cannot detect composite types like pointer, channel, and so on.
// If you know the type is comparable, please use `NewComparableMap()` instead.
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

// Create new map from existing map
//
// # Arguments
// * `m`: `map[K]V` - Existing map
//
// # Returns
// * `IMap[K, V]` - Created map
//
// # Notes
// For some predefined primitive types, we detect and create `ComparableMap[K, V]`.
// This mechanism cannot detect composite types like pointer, channel, and so on.
// If you know the type is comparable, please use `NewComparableMap()` instead.
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

// Create new comparable map
//
// # Returns
// * `IComparable[K, V]` - Created comparable map
func NewComparableMap[K, V comparable]() IComparableMap[K, V] {
	return &ComparableMap[K, V]{
		Map[K,V]{
			item: make(map[K]V, 0),
		},
	}
}

// Create new comparable map from existing map
//
// # Arguments
// * `m`: `map[K]V` - Existing comparable map
//
// # Returns
// * `IComparable[K, V]` - Created comparable map
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

// Free function of equality check
//
// # Arguments
// * `m1`, `m2`: `IComparableMap[K, V]` - Maps to be compared
//
// # Returns
// * `bool` - Whether they are equal
func Equal[K, V comparable](m1, m2 IComparableMap[K, V]) bool {
	return m1.Equal(m2)
}

// Free function of equality check
//
// # Arguments
// * `m1`: `IComparableMap[K, V1]` - Maps to be compared
// * `m2`: `IComparableMap[K, V2]` - Maps to be compared
// * `eq`: `func (V1, V2) bool`
//
// # Returns
// * `bool` - Whether they are equal
func EqualFunc[K comparable, V1, V2 any](m1 IMap[K, V1], m2 IMap[K, V2], eq func (V1, V2) bool) bool {
	if m1.Size() != m2.Size() {
		return false
	}

	for k := range m1.Keys() {
		// # Issue
		// Go 1.19 compiler cannot estimate type correctly, and fail to build;
		// `cannot use k (variable of type int) as type K in argument to m1.Get` etc.
		// As a workaround, we dynamically assert type
		key := any(k).(K)
		v1, _ := m1.Get(key)
		if v2, ok := m2.Get(key); (!ok) || (!eq(v1, v2)) {
			return false
		}
	}

	return true
}
