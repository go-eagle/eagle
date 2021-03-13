package cache

import (
	"fmt"
	"testing"
)

func TestLRU(t *testing.T) {
	lru := NewLRU(3)
	lru.Set(1, 1)
	lru.Set(2, 2)
	lru.Set(3, 3)

	// Least 1 -> 2 -> 3 Most
	lru.ShowQueue()

	// 2
	fmt.Println(lru.Get(2))
	// Least 1 -> 3 -> 2 Most
	lru.ShowQueue()

	lru.Set(1, 100)
	// Least 3 -> 2 -> 1 Most
	lru.ShowQueue()

	lru.Set(4, 4)
	// Least 2 -> 1 -> 4 Most
	lru.ShowQueue()
}
