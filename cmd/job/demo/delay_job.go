package demo

import (
	"log"
	"time"
)

type DelayJob struct {
	count int
}

func (d *DelayJob) Run() {
	time.Sleep(2 * time.Second)
	d.count++
	log.Printf("%d: hello world\n", d.count)
}
