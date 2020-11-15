// 异步缓存
// 场景1：本地缓存更新，可以使用cron同步远程数据到本地
package cache

import (
	"fmt"
	"sync"
	"testing"

	"github.com/robfig/cron"
)

type store struct {
	sync.RWMutex
	cron *cron.Cron
	data map[int]string
}

var storeObj *store

// go test -v sync_cache_example_test.go
func TestSyncCache(t *testing.T) {
	storeObj = &store{
		data: make(map[int]string, 0),
		cron: cron.New(),
	}
	storeObj.syncStore()
	storeObj.cron.AddFunc("*/10 * * * * *", storeObj.syncStore)
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

func (s *store) syncStore() {
	dataTmp := getRemoteData()
	s.Lock()
	s.data = dataTmp
	s.Unlock()
}

func (s *store) Get(id int) string {
	s.RLock()
	defer s.RUnlock()
	resp, _ := s.data[id]
	return resp
}

func getRemoteData() map[int]string {
	resp := make(map[int]string)
	resp[1] = "test1"
	resp[2] = "test2"
	resp[3] = "test3"
	resp[4] = "test4"
	resp[5] = "test5"
	return resp
}
