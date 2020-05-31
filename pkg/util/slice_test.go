package util

import (
	"reflect"
	"testing"
)

func TestStringSliceReflectEqual(t *testing.T) {
	cases := []struct {
		in, want []string
	}{
		{[]string{"q", "w", "e", "r", "t"}, []string{"q", "w", "a", "s", "z", "x"}},
	}
	for _, c := range cases {
		result := StringSliceReflectEqual(c.in, c.want)
		if !result {
			t.Errorf("StringSliceReflectEqual(%q) == %q", c.in, c.want)
		}
	}
}

func BenchmarkDeepEqual(b *testing.B) {
	sa := []string{"q", "w", "e", "r", "t"}
	sb := []string{"q", "w", "a", "s", "z", "x"}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		StringSliceReflectEqual(sa, sb)
	}
}

func TestStringSliceEqual(t *testing.T) {
	cases := []struct {
		in, want []string
	}{
		{[]string{"q", "w", "e", "r", "t"}, []string{"q", "w", "a", "s", "z", "x"}},
	}
	for _, c := range cases {
		result := StringSliceEqual(c.in, c.want)
		if !result {
			t.Errorf("StringSliceEqual(%q) == %q", c.in, c.want)
		}
	}
}

func BenchmarkEqual(b *testing.B) {
	sa := []string{"q", "w", "e", "r", "t"}
	sb := []string{"q", "w", "a", "s", "z", "x"}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		StringSliceEqual(sa, sb)
	}
}

func TestUint64SliceReverse(t *testing.T) {
	cases := []struct {
		in, want []uint64
	}{
		{[]uint64{1, 2, 3, 4, 5}, []uint64{5, 4, 3, 2, 1}},
	}
	for _, c := range cases {
		got := Uint64SliceReverse(c.in)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("SliceReverseUint64(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}

func TestDeleteElemFromUint64Slice(t *testing.T) {
	type args struct {
		i int
		a []uint64
	}
	tests := []struct {
		name string
		args args
		want []uint64
	}{
		{"delete 3 from slice", args{
			i: 2,
			a: []uint64{1, 2, 3, 4, 5},
		}, []uint64{1, 2, 5, 4}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DeleteElemFromUint64Slice(tt.args.i, tt.args.a); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteElemFromUint64Slice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteElemOrderFromUint64Slice(t *testing.T) {
	type args struct {
		i int
		a []uint64
	}
	tests := []struct {
		name string
		args args
		want []uint64
	}{
		{"delete middle element 3 from slice", args{
			i: 2,
			a: []uint64{1, 2, 3, 4, 5},
		}, []uint64{1, 2, 4, 5}},
		{"delete first element 1 from slice", args{
			i: 0,
			a: []uint64{1, 2, 3, 4, 5},
		}, []uint64{2, 3, 4, 5}},
		{"delete last element 5 from slice", args{
			i: 4,
			a: []uint64{1, 2, 3, 4, 5},
		}, []uint64{1, 2, 3, 4}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DeleteElemOrderFromUint64Slice(tt.args.i, tt.args.a); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteElemOrderFromUint64Slice() = %v, want %v", got, tt.want)
			}
		})
	}
}
