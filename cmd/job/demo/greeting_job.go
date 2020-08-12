package demo

import "fmt"

type GreetingJob struct {
	Name string
}

func (g GreetingJob) Run() {
	fmt.Println("Hello ", g.Name)
}
