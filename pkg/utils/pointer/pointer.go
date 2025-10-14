package pointer

import "reflect"

// Ptr returns a pointer to its argument.
// For example, Ptr[int](100) returns *int pointing to 100.
func Ptr[T any](t T) *T { return &t }

// Value returns the value that the pointer p points to.
// If p is nil, it returns the zero value of the type T.
func Value[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}

// IsStructPtr checks if the given interface is a pointer to a struct.
func IsStructPtr(x any) bool {
	if x == nil {
		return false
	}
	t := reflect.TypeOf(x)
	return t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct
}
