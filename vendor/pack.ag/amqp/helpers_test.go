package amqp

import (
	"reflect"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func testEqual(x, y interface{}) bool {
	return cmp.Equal(x, y, compareOpts(x, y)...)
}

func testDiff(x, y interface{}) string {
	return cmp.Diff(x, y, compareOpts(x, y)...)
}

func compareOpts(x, y interface{}) []cmp.Option {
	return cmp.Options{
		deepAllowUnexported(x, y),
		cmpopts.EquateNaNs(),
	}
}

// from https://github.com/google/go-cmp/issues/40
func deepAllowUnexported(vs ...interface{}) cmp.Option {
	m := make(map[reflect.Type]struct{})
	for _, v := range vs {
		structTypes(reflect.ValueOf(v), m)
	}
	var types []interface{}
	for t := range m {
		types = append(types, reflect.New(t).Elem().Interface())
	}
	return cmp.AllowUnexported(types...)
}

func structTypes(v reflect.Value, m map[reflect.Type]struct{}) {
	if !v.IsValid() {
		return
	}
	switch v.Kind() {
	case reflect.Ptr:
		if !v.IsNil() {
			structTypes(v.Elem(), m)
		}
	case reflect.Interface:
		if !v.IsNil() {
			structTypes(v.Elem(), m)
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			structTypes(v.Index(i), m)
		}
	case reflect.Map:
		for _, k := range v.MapKeys() {
			structTypes(v.MapIndex(k), m)
		}
	case reflect.Struct:
		m[v.Type()] = struct{}{}
		for i := 0; i < v.NumField(); i++ {
			structTypes(v.Field(i), m)
		}
	}
}
