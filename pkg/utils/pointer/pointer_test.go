package pointer

import (
	"reflect"
	"testing"
)

func TestPtr(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		val := 100
		ptr := Ptr(val)

		if ptr == nil {
			t.Fatal("ToPtr() returned nil")
		}
		if *ptr != val {
			t.Errorf("expected value %d, got %d", val, *ptr)
		}
		if reflect.TypeOf(ptr).Kind() != reflect.Ptr {
			t.Errorf("expected a pointer type, got %T", ptr)
		}
	})

	t.Run("string", func(t *testing.T) {
		val := "hello"
		ptr := Ptr(val)

		if ptr == nil {
			t.Fatal("Ptr() returned nil")
		}
		if *ptr != val {
			t.Errorf("expected value %q, got %q", val, *ptr)
		}
	})
}

func TestIsStructPtr(t *testing.T) {
	type testStruct struct {
		Name string
	}
	var nilStructPtr *testStruct

	testCases := []struct {
		name     string
		input    any
		expected bool
	}{
		{"pointer to struct", &testStruct{Name: "test"}, true},
		{"nil pointer to struct", nilStructPtr, true},
		{"struct value", testStruct{Name: "test"}, false},
		{"pointer to int", Ptr(123), false},
		{"int value", 123, false},
		{"pointer to string", Ptr("hello"), false},
		{"string value", "hello", false},
		{"nil interface", nil, false},
		{"pointer to pointer to struct", Ptr(&testStruct{}), false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := IsStructPtr(tc.input)
			if result != tc.expected {
				t.Errorf("IsStructPtr(%v) = %v; want %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestValue(t *testing.T) {
	type testStruct struct {
		Name string
	}

	testCases := []struct {
		name     string
		input    any
		expected any
	}{
		{"non-nil int pointer", Ptr(123), 123},
		{"nil int pointer", (*int)(nil), 0},
		{"non-nil string pointer", Ptr("hello"), "hello"},
		{"nil string pointer", (*string)(nil), ""},
		{"non-nil struct pointer", &testStruct{Name: "test"}, testStruct{Name: "test"}},
		{"nil struct pointer", (*testStruct)(nil), testStruct{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var result any
			switch v := tc.input.(type) {
			case *int:
				result = Value(v)
			case *string:
				result = Value(v)
			case *testStruct:
				result = Value(v)
			default:
				t.Fatalf("unhandled test case type: %T", tc.input)
			}

			if result != tc.expected {
				t.Errorf("Value() = %v; want %v", result, tc.expected)
			}
		})
	}
}
