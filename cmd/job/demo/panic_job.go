package demo

import "fmt"

type PanicJob struct {
	count int
}

func (p *PanicJob) Run() {
	p.count++
	if p.count == 1 {
		panic("oooooooooooooops!!!")
	}

	fmt.Println("hello world")
}
