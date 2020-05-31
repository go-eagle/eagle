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

// Uint64ShuffleSlice 对slice进行随机
func Uint64ShuffleSlice(a []uint64) []uint64 {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(a), func(i, j int) {
		a[i], a[j] = a[j], a[i]
	})
	return a
}

// see: https://yourbasic.org/golang/

// DeleteElemFromUint64Slice 从slice删除元素 fast version, 会改变顺序
// i：slice的索引值
// a: slice
func DeleteElemFromUint64Slice(i int, a []uint64) []uint64 {
	// Remove the element at index i from a.
	a[i] = a[len(a)-1] // Copy last element to index i.
	a[len(a)-1] = 0    // Erase last element (write zero value).
	a = a[:len(a)-1]   // Truncate slice.

	return a
}

// DeleteElemOrderFromUint64Slice 从slice删除元素 slow version, 保持原有顺序
// i：slice的索引值
// a: slice
func DeleteElemOrderFromUint64Slice(i int, a []uint64) []uint64 {
	// Remove the element at index i from a.
	copy(a[i:], a[i+1:]) // Shift a[i+1:] left one index.
	a[len(a)-1] = 0      // Erase last element (write zero value).
	a = a[:len(a)-1]     // Truncate slice.

	return a
}
