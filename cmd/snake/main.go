package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/1024casts/snake/cmd/snake/internal/project"
	"github.com/1024casts/snake/cmd/snake/internal/upgrade"
)

var (
	// Version is the version of the compiled software.
	Version string = "v0.2.0"

	rootCmd = &cobra.Command{
		Use:     "snake",
		Short:   "Snake: An elegant toolkit for Go microservices.",
		Long:    `Snake: An elegant toolkit for Go microservices.`,
		Version: Version,
	}
)

func init() {
	rootCmd.AddCommand(project.CmdNew)
	rootCmd.AddCommand(upgrade.CmdUpgrade)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
