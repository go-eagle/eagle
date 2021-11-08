package main

import (
	"fmt"
	"log"

	"github.com/go-eagle/eagle/pkg/config"
)

func main() {
	type Service struct {
		Name    string
		Version string
	}

	c := config.New(".")

	var svc Service
	err := c.Load("config", &svc)
	if err != nil {
		log.Fatalf("load config err: %+v", err)
	}
	fmt.Println("service name: ", svc.Name)

	conf, err := c.LoadWithType("server", "yaml")
	if err != nil {
		log.Fatalf("load server err: %+v", err)
	}
	fmt.Println("http addr: ", conf.GetString("http.addr"))
	fmt.Println("http port: ", conf.GetInt("http.port"))
}
