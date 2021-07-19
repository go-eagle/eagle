package main

import (
	"log"

	"github.com/1024casts/snake/cmd/snake/internal/run"

	"github.com/1024casts/snake/cmd/snake/internal/cache"

	"github.com/spf13/cobra"

	"github.com/1024casts/snake/cmd/snake/internal/project"
	"github.com/1024casts/snake/cmd/snake/internal/upgrade"
)

var (
	// Version is the version of the compiled software.
	Version = "v0.3.0"

	rootCmd = &cobra.Command{
		Use:     "snake",
		Short:   "Snake: An develop kit for Go microservices.",
		Long:    `Snake: An develop kit for Go microservices.`,
		Version: Version,
	}
)

func init() {
	rootCmd.AddCommand(project.CmdNew)
	rootCmd.AddCommand(run.CmdRun)
	rootCmd.AddCommand(cache.CmdCache)
	rootCmd.AddCommand(upgrade.CmdUpgrade)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
