package async

import (
	"context"

	"github.com/go-eagle/eagle/pkg/log"
	"github.com/go-eagle/eagle/pkg/utils"
)

// Go 异步执行 asyncFunc() 函数，会进行 recover() 操作，如果出现 panic() 则会记录日志
// name 用于 recover 后的日志和 metrics 统计
func Go(ctx context.Context, name string, asyncFunc func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				stack := utils.StackTrace(name, err)
				log.WithContext(ctx).Errorf("async: name: %s panic: %s stack: %s", name, err, stack)
				// TODO: metrics
			}
		}()

		asyncFunc()
	}()
}
