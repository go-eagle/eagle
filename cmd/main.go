package main

import (
	"os"

	"github.com/urfave/cli"

	"github.com/1024casts/snake/cmd/snake/new"
)

const Version = "1.0.0"

func main() {
	app := cli.NewApp()
	app.Name = "snake"
	app.Usage = "snake tools"
	app.Version = Version
	app.Commands = []cli.Command{
		new.Cmd,
	}
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
