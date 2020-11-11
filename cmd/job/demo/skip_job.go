package demo

import (
	"log"
	"sync/atomic"
	"time"
)

// SkipJob define a struct
type SkipJob struct {
	count int32
}

// Run run job
func (d *SkipJob) Run() {
	atomic.AddInt32(&d.count, 1)
	log.Printf("%d: hello world\n", d.count)
	if atomic.LoadInt32(&d.count) == 1 {
		time.Sleep(2 * time.Second)
	}
}
