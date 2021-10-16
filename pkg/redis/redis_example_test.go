// 并发获取数据
package redis

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// go test -v redis_example_test.go
func TestOneRedisData(t *testing.T) {
	t1 := time.Now()
	for i := 0; i < 120; i++ {
		getRemoteOneRedisData(i)
	}
	fmt.Println("Test_OneRedisData cost: ", time.Since(t1))
}

func TestPipelineRedisData(t *testing.T) {
	t1 := time.Now()
	ids := make([]int, 0, 120)
	for i := 0; i < 120; i++ {
		ids = append(ids, i)
	}
	getRemotePipelineRedisData(ids)
	fmt.Println("Test_PipelineRedisData cost: ", time.Since(t1))
}

func TestGoroutinePipelineRedisData(t *testing.T) {
	t1 := time.Now()
	ids := make([]int, 0, 120)
	for i := 0; i < 120; i++ {
		ids = append(ids, i)
	}
	getGoroutinePipelineRedisData(ids)
	fmt.Println("Test_GoroutinePipelineRedisData cost: ", time.Since(t1))
}

func getRemoteOneRedisData(i int) int {
	// 模拟单个redis请求，定义为600us
	time.Sleep(600 * time.Microsecond)
	return i
}

// getRemotePipelineRedisData 多id批量获取
func getRemotePipelineRedisData(i []int) []int {
	length := len(i)
	// 使用pipeline的情况下，单个redis数据，为500us
	time.Sleep(time.Duration(length) * 500 * time.Microsecond)
	return i
}

// getGoroutinePipelineRedisData 分段后通过goroutine并发获取数据
func getGoroutinePipelineRedisData(ids []int) []int {
	idsNew := make(map[int][]int)
	idsNew[0] = ids[0:30]
	idsNew[1] = ids[30:60]
	idsNew[2] = ids[60:90]
	idsNew[3] = ids[90:120]

	resp := make([]int, 0, 120)
	var wg sync.WaitGroup
	for j := 0; j < 4; j++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, j int) {
			resp = append(resp, getRemotePipelineRedisData(idsNew[j])...)
			wg.Done()
		}(&wg, j)
	}
	wg.Wait()

	return resp
}
