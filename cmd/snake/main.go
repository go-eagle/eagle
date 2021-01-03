package main

import (
	"os"

	"github.com/1024casts/snake/cmd/snake/new"
	"github.com/urfave/cli"
)

const Version = "0.1.0"

func main() {
	app := cli.NewApp()
	app.Name = "snake"
	app.Usage = "snake tools"
	app.Version = Version
	app.Commands = []cli.Command{
		{
			Name:            "new",
			Aliases:         []string{"n"},
			Usage:           "Create Snake template project",
			Action:          new.CreateProject,
			SkipFlagParsing: false,
			UsageText:       new.NewProjectHelpTemplate,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "d",
					Value:       "",
					Usage:       "Specify the directory of the project",
					Destination: &new.Project.Path,
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
