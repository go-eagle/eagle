package job

import (
	"time"

	"github.com/robfig/cron/v3"

	"github.com/go-eagle/eagle/cmd/job/example"
	"github.com/go-eagle/eagle/pkg/log"
)

// Run 计划任务
// see: https://mp.weixin.qq.com/s/Ak7RBv1NuS-VBeDNo8_fww
//
// cron 内置3个用得比较多的JobWrapper：
//
// Recover：捕获内部Job产生的 panic；
// DelayIfStillRunning：触发时，如果上一次任务还未执行完成（耗时太长），则等待上一次任务完成之后再执行；
// SkipIfStillRunning：触发时，如果上一次任务还未完成，则跳过此次执行。
func Run() {
	c := cron.New()
	// support to second
	// c = cron.New(cron.WithSeconds())
	// demo
	_, err := c.AddFunc("* */5 * * *", func() {
		log.Infof("test cron, time: %d ", time.Now().Unix())
	})
	if err != nil {
		log.Warnf("cron AddFunc err, %+v", err)
		return
	}

	// test recover
	_, _ = c.AddJob("@every 1s", cron.NewChain(cron.Recover(cron.DefaultLogger)).Then(&example.PanicJob{}))

	// test DelayIfStillRunning
	_, _ = c.AddJob("@every 1s", cron.NewChain(cron.DelayIfStillRunning(cron.DefaultLogger)).Then(&example.DelayJob{}))

	// test SkipIfStillRunning
	_, _ = c.AddJob("@every 1s", cron.NewChain(cron.SkipIfStillRunning(cron.DefaultLogger)).Then(&example.SkipJob{}))

	// 执行具体的任务
	// _, _ = c.AddJob("@every 3s", example.GreetingJob{"dj"})

	c.Start()
}
