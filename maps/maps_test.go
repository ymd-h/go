package maps

import (
	"reflect"
	"testing"

	y "github.com/ymd-h/go/testing"
)

func TestNewMap(t *testing.T) {
	var i int
	tests := []struct{
		name string
		args any
		want any
	}{
		{
			name: "bool",
			args: true,
			want: NewComparableMap[int, bool](),
		},
		{
			name: "string",
			args: "",
			want: NewComparableMap[int, string](),
		},
		{
			name: "int",
			args: int(1),
			want: NewComparableMap[int, int](),
		},
		{
			name: "uint",
			args: uint(1),
			want: NewComparableMap[int, uint](),
		},
		{
			name: "*int (undetectable)",
			args: &i,
			want: &Map[int, *int]{ item: map[int]*int{} },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(*testing.T) {
			wantName := reflect.TypeOf(tt.want).Name()
			var got string

			switch any(tt.args).(type) {
			case bool:
				got = reflect.TypeOf(NewMap[int, bool]()).Name()
			case string:
				got = reflect.TypeOf(NewMap[int, string]()).Name()
			case int:
				got = reflect.TypeOf(NewMap[int, int]()).Name()
			case uint:
				got = reflect.TypeOf(NewMap[int, uint]()).Name()
			case *int:
				got = reflect.TypeOf(NewMap[int, *int]()).Name()
			}

			if wantName != got {
				t.Errorf("%v != %v\n", got, wantName)
			}
		})
	}
}


func TestNewMapFrom(t *testing.T) {
	tests := []struct {
		name string
		args any
		want any
	}{
		{
			name: "bool",
			args: map[int]bool{},
			want: NewMap[int, bool](),
		},
		{
			name: "string",
			args: map[int]string{},
			want: NewMap[int, string](),
		},
		{
			name: "int",
			args: map[int]int{},
			want: NewMap[int, int](),
		},
		{
			name: "uint",
			args: map[int]uint{},
			want: NewMap[int, uint](),
		},
		{
			name: "*int (unsupported)",
			args: map[int] *int{},
			want: NewMap[int, *int](),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(*testing.T) {
			wantName := reflect.TypeOf(tt.want).Name()
			var got string

			switch m := tt.args.(type) {
			case map[int]bool:
				got = reflect.TypeOf(NewMapFrom(m)).Name()
			case map[int]string:
				got = reflect.TypeOf(NewMapFrom(m)).Name()
			case map[int]int:
				got = reflect.TypeOf(NewMapFrom(m)).Name()
			case map[int]uint:
				got = reflect.TypeOf(NewMapFrom(m)).Name()
			case map[int]*int:
				got = reflect.TypeOf(NewMapFrom(m)).Name()
			}

			if wantName != got {
				t.Errorf("%v != %v\n", got, wantName)
			}
		})
	}
}


func TestGet(t *testing.T) {
	tests := []struct {
		name string
		args map[int]int
		key int
		want int
		wantOk bool
	}{
		{
			name: "empty",
			args: map[int]int{},
			key: 1,
			want: 0,
			wantOk: false,
		},
		{
			name: "simple",
			args: func() map[int]int {
				m := map[int]int{}
				m[0] = 1
				return m
			}(),
			key: 0,
			want: 1,
			wantOk: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(*testing.T) {
			m := NewMapFrom(tt.args)
			v, ok := m.Get(tt.key)

			if ok != tt.wantOk {
				t.Errorf("%v != %v\n", ok, tt.wantOk)
			} else {
				if ok && v != tt.want {
					t.Errorf("%v != %v\n", v, tt.want)
				}
			}
		})
	}
}


func TestSet(t *testing.T) {
	type test struct {
		init IMap[string, int]
		key string
		value int
		wantBefore int
		wantBeforeOk bool
		wantAfter int
	}

	y.NewTest[test](t).
		Add("simple", test{
			init: NewMap[string, int](),
			key: "a",
			value: 1,
			wantBefore: 0,
			wantBeforeOk: false,
			wantAfter: 1,
		}).
		Add("overwrite", test{
			init: func() IMap[string, int] {
				m := NewMap[string, int]()
				m.Set("b", 2)
				return m
			}(),
			key: "b",
			value: 10,
			wantBefore: 2,
			wantBeforeOk: true,
			wantAfter: 10,
		}).
		Run(func(tt *testing.T, data test) {
			m := data.init
			before, beforeOk := m.Get(data.key)
			y.AssertEqual(t, beforeOk, data.wantBeforeOk)

			if beforeOk {
				y.AssertEqual(t, before, data.wantBefore)
			}

			m.Set(data.key, data.value)
			after, afterOk := m.Get(data.key)
			y.AssertEqual(t, afterOk, true)
			y.AssertEqual(t, after, data.wantAfter)
	})
}
