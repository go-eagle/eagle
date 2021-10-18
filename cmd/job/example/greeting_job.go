package example

import "fmt"

// GreetingJob define struct
type GreetingJob struct {
	Name string
}

// Run run job
func (g GreetingJob) Run() {
	fmt.Println("[greeting] Hello ", g.Name)
}
