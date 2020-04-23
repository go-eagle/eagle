package schedule

import (
	"time"

	"github.com/lexkong/log"
	"github.com/robfig/cron/v3"
)

// Init 初始化计划任务
func Init() {
	c := cron.New()
	// demo
	_, err := c.AddFunc("*/3 * * * *", func() {
		log.Infof("test cron, time: %d ", time.Now().Unix())
	})
	if err != nil {
		log.Warnf("cron AddFunc err, %+v", err)
		return
	}

	c.Start()
}
