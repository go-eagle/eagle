package demo

import (
	"log"
	"sync/atomic"
	"time"
)

type SkipJob struct {
	count int32
}

func (d *SkipJob) Run() {
	atomic.AddInt32(&d.count, 1)
	log.Printf("%d: hello world\n", d.count)
	if atomic.LoadInt32(&d.count) == 1 {
		time.Sleep(2 * time.Second)
	}
}
