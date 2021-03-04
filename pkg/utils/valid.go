package utils

import (
	"reflect"
)

// IsZero 检查是否是零值
func IsZero(i ...interface{}) bool {
	ret := false
	for _, j := range i {
		v := reflect.ValueOf(j)
		if isZero(v) {
			return true
		}
	}
	return ret
}

func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Slice:
		return v.IsNil()
	case reflect.Invalid:
		return true
	default:
		z := reflect.Zero(v.Type())
		return reflect.DeepEqual(z.Interface(), v.Interface())
	}
}
