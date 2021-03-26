// 异步缓存
// 场景1：本地缓存更新，可以使用cron同步远程数据到本地
package cache

import (
	"fmt"
	"sync"
	"testing"
)

// go test -v sync_cache_test.go
func TestSyncCache(t *testing.T) {
	storeObj = NewSyncStore()

	storeObj.syncStore(getRemoteData)
	storeObj.cron.AddFunc("*/1 * * * * *", func() {
		storeObj.syncStore(getRemoteData)
	})
	storeObj.cron.Start()

	var wg sync.WaitGroup
	for i := 1; i < 6; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, i int) {
			fmt.Println(storeObj.Get(i))
			wg.Done()
		}(&wg, i)
	}
	wg.Wait()
}

func getRemoteData() map[int]interface{} {
	resp := make(map[int]interface{})
	resp[1] = "test1"
	resp[2] = "test2"
	resp[3] = "test3"
	resp[4] = "test4"
	resp[5] = "test5"
	return resp
}
