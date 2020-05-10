package schedule

import (
	"time"

	"github.com/robfig/cron/v3"

	"github.com/1024casts/snake/pkg/log"
)

// Init 初始化计划任务
func Init() {
	c := cron.New()
	// demo
	_, err := c.AddFunc("* */5 * * *", func() {
		log.Infof("test cron, time: %d ", time.Now().Unix())
	})
	if err != nil {
		log.Warnf("cron AddFunc err, %+v", err)
		return
	}

	c.Start()
}
