package demo

import "fmt"

// PanicJob define a struct
type PanicJob struct {
	count int
}

// Run run job
func (p *PanicJob) Run() {
	p.count++
	if p.count == 1 {
		panic("oooooooooooooops!!!")
	}

	fmt.Println("hello world")
}
