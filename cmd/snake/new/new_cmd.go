package new

import "github.com/urfave/cli"

var Cmd = cli.Command{
	Name:            "new",
	Aliases:         []string{"n"},
	Usage:           "Create Snake template project",
	Action:          CreateProject,
	SkipFlagParsing: false,
	UsageText:       NewProjectHelpTemplate,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "d",
			Value:       "",
			Usage:       "Specify the directory of the project",
			Destination: &project.Path,
		},
	},
}
