package example

import (
	"log"
	"time"
)

// DelayJob delay job
type DelayJob struct {
	count int
}

// Run run job
func (d *DelayJob) Run() {
	time.Sleep(2 * time.Second)
	d.count++
	log.Printf("[deplay] %d: hello world\n", d.count)
}
