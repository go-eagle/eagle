package main

import (
	"github.com/go-eagle/eagle/pkg/log"
)

type GreetingJob struct {
	Name string
}

func (g GreetingJob) Run() {
	log.Info("Hello ", g.Name)
}

type SendEmail struct {
	Name string
}

func (g SendEmail) Run() {
	log.Info("Send mail to... ", g.Name)
}
