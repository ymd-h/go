package maps

import (
	"reflect"
	"testing"
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
