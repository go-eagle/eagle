package main

import (
	"github.com/go-eagle/eagle/pkg/config"
	"github.com/go-eagle/eagle/pkg/log"
)

func main() {
	// load config
	_ = config.New(".")

	// print log using default filename
	log.Init()
	log.Info("test log")

	// print log using custom filename
	log.Init(log.WithFilename("custom-filename"))
	log.Info("test log with filename")

	// print log using custom dir and filename
	log.Init(log.WithFilename("user/custom"))
	log.Info("test log with dir and filename")
}
