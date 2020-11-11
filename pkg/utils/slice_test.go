package utils

import (
	"reflect"
	"strconv"
	"sync"
	"testing"
)

func TestStringSliceReflectEqual(t *testing.T) {
	cases := []struct {
		name     string
		inA, inB []string
		want     bool
	}{
		{"test slice is not equal", []string{"q", "w", "e", "r", "t"}, []string{"q", "w", "a", "s", "z", "x"}, false},
		{"test slice is equal", []string{"q", "w", "e", "r", "t"}, []string{"q", "w", "e", "r", "t"}, true},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringSliceReflectEqual(tt.inA, tt.inB); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteElemOrderFromUint64Slice() = %v, want %v", got, tt.want)
			}
		})
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
	type args struct {
		a []string
		b []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"test slice is equal", args{
			a: []string{"golang", "ruby", "java", "rust"},
			b: []string{"golang", "ruby", "java", "rust"},
		}, true},
		{"test slice is not equal", args{
			a: []string{"golang", "ruby", "java", "rust"},
			b: []string{"php", "ruby", "java", "rust"},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringSliceEqual(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("StringSliceEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceShuffle(t *testing.T) {
	slice := []interface{}{"a", "b", "c", "d", "e", "f"}
	SliceShuffle(slice)
	t.Log(slice)
	for _, v := range slice {
		t.Log(v.(string))
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

func TestUint64DeleteElemInSlice(t *testing.T) {
	type args struct {
		i int
		a []uint64
	}
	tests := []struct {
		name string
		args args
		want []uint64
	}{
		{"delete middle 3 from slice", args{
			i: 2,
			a: []uint64{1, 2, 3, 4, 5},
		}, []uint64{1, 2, 5, 4}},
		{"delete first 1 from slice", args{
			i: 0,
			a: []uint64{1, 2, 3, 4, 5},
		}, []uint64{5, 2, 3, 4}},
		{"delete last 5 from slice", args{
			i: 4,
			a: []uint64{1, 2, 3, 4, 5},
		}, []uint64{1, 2, 3, 4}},
		{"delete element out of last range slice", args{
			i: 5,
			a: []uint64{1, 2, 3, 4, 5},
		}, []uint64{1, 2, 3, 4, 5}},
		{"delete element out of first range slice", args{
			i: -1,
			a: []uint64{1, 2, 3, 4, 5},
		}, []uint64{1, 2, 3, 4, 5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Uint64DeleteElemInSlice(tt.args.i, tt.args.a); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Uint64DeleteElemInSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint64DeleteElemInSliceWithOrder(t *testing.T) {
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
			if got := Uint64DeleteElemInSliceWithOrder(tt.args.i, tt.args.a); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Uint64DeleteElemInSliceWithOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint64ShuffleSlice(t *testing.T) {
	type args struct {
		a []uint64
	}
	tests := []struct {
		name string
		args args
		want []uint64
	}{
		{"test gen rand slice", args{a: []uint64{1, 2, 3, 4, 5}}, []uint64{5, 1, 3, 4, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Uint64ShuffleSlice(tt.args.a); reflect.DeepEqual(got, tt.want) {
				t.Errorf("Uint64ShuffleSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringSliceContains(t *testing.T) {
	type args struct {
		ss []string
		s  string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"test contain string in slice", args{
			ss: []string{"golang", "java", "rust", "python", "php"},
			s:  "golang",
		}, true},
		{"test not contain string in slice", args{
			ss: []string{"golang", "java", "rust", "python", "php"},
			s:  "ruby",
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringSliceContains(tt.args.ss, tt.args.s); got != tt.want {
				t.Errorf("StringSliceContains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsInSlice(t *testing.T) {
	type args struct {
		value interface{}
		sli   interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"test int in slice", args{
			value: 1,
			sli:   []int{1, 2, 3},
		}, true},
		{"test int not in slice", args{
			value: 4,
			sli:   []int{1, 2, 3},
		}, false},
		{"test string in slice", args{
			value: "golang",
			sli:   []string{"golang", "mysql", "redis"},
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInSlice(tt.args.value, tt.args.sli); got != tt.want {
				t.Errorf("IsInSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

// go test -bench="."
// benchamark slice compare
type Something struct {
	roomID   int
	roomName string
}

// BenchmarkDefaultSlice default slice
func BenchmarkDefaultSlice(b *testing.B) {
	b.ReportAllocs()
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			for i := 0; i < 120; i++ {
				output := make([]Something, 0)
				output = append(output, Something{
					roomID:   i,
					roomName: strconv.Itoa(i),
				})
			}
			wg.Done()
		}(&wg)
	}
	wg.Wait()
}

// BenchmarkPreAllocSlice 预分配
func BenchmarkPreAllocSlice(b *testing.B) {
	b.ReportAllocs()
	var wg sync.WaitGroup

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			output := make([]Something, 0, 120)
			for i := 0; i < 120; i++ {
				output = append(output, Something{
					roomID:   i,
					roomName: strconv.Itoa(i),
				})
			}
			wg.Done()
		}(&wg)
	}
	wg.Wait()
}

// BenchmarkSyncPoolSlice 使用 sync pool
func BenchmarkSyncPoolSlice(b *testing.B) {
	b.ReportAllocs()
	var wg sync.WaitGroup
	var SomethingPool = sync.Pool{
		New: func() interface{} {
			b := make([]Something, 120)
			return &b
		},
	}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			obj := SomethingPool.Get().(*[]Something)
			for i := 0; i < 120; i++ {
				some := *obj
				some[i].roomID = i
				some[i].roomName = strconv.Itoa(i)
			}
			SomethingPool.Put(obj)
			wg.Done()
		}(&wg)
	}
	wg.Wait()
}
