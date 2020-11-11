package lock

import (
	"fmt"
	"math/rand"
	"time"
)

// ExampleTry test try lock
func ExampleTry() {
	var mu Mutex
	go func() { // 启动一个goroutine持有一段时间的锁
		mu.Lock()
		time.Sleep(time.Duration(rand.Intn(2)) * time.Second)
		mu.Unlock()
	}()

	time.Sleep(time.Second)

	ok := mu.TryLock() // 尝试获取到锁
	if ok {            // 获取成功
		fmt.Println("got the lock")
		// do something
		mu.Unlock()
		return
	}

	// 没有获取到
	fmt.Println("can't get the lock")
}
