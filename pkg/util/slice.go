package util

import (
	"math/rand"
	"reflect"
	"time"
)

// StringSliceReflectEqual 判断 string和slice 是否相等
// 因为使用了反射，所以效率较低，可以看benchmark结果
func StringSliceReflectEqual(a, b []string) bool {
	return reflect.DeepEqual(a, b)
}

// StringSliceEqual 判断 string和slice 是否相等
// 使用了传统的遍历方式
func StringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	// reflect.DeepEqual的结果保持一致
	if (a == nil) != (b == nil) {
		return false
	}

	// bounds check 边界检查
	// 避免越界
	b = b[:len(a)]
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}

// Uint64SliceReverse 对uint64 slice 反转
func Uint64SliceReverse(a []uint64) []uint64 {
	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}

	return a
}

// StringSliceContains 字符串切片中是否包含另一个字符串
// 来自go源码 net/http/server.go
func StringSliceContains(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}

// IsInSlice 判断某一值是否在slice中
// 因为使用了反射，所以时间开销比较大，使用中根据实际情况进行选择
func IsInSlice(value interface{}, sli interface{}) bool {
	switch reflect.TypeOf(sli).Kind() {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(sli)
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(value, s.Index(i).Interface()) {
				return true
			}
		}
	}
	return false
}

// Uint64ShuffleSlice 对slice进行随机
func Uint64ShuffleSlice(a []uint64) []uint64 {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(a), func(i, j int) {
		a[i], a[j] = a[j], a[i]
	})
	return a
}

// see: https://yourbasic.org/golang/

// Uint64DeleteElemInSlice 从slice删除元素
// fast version, 会改变顺序
// i：slice的索引值
// s: slice
func Uint64DeleteElemInSlice(i int, s []uint64) []uint64 {
	if i < 0 || i > len(s)-1 {
		return s
	}
	// Remove the element at index i from s.
	s[i] = s[len(s)-1] // Copy last element to index i.
	s[len(s)-1] = 0    // Erase last element (write zero value).
	s = s[:len(s)-1]   // Truncate slice.

	return s
}

// Uint64DeleteElemInSliceWithOrder 从slice删除元素
// slow version, 保持原有顺序
// i：slice的索引值
// s: slice
func Uint64DeleteElemInSliceWithOrder(i int, s []uint64) []uint64 {
	if i < 0 || i > len(s)-1 {
		return s
	}
	// Remove the element at index i from s.
	copy(s[i:], s[i+1:]) // Shift s[i+1:] left one index.
	s[len(s)-1] = 0      // Erase last element (write zero value).
	s = s[:len(s)-1]     // Truncate slice.

	return s
}
