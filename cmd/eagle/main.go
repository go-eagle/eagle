package main

import (
	"log"

	"github.com/go-eagle/eagle/cmd/eagle/internal/proto"

	"github.com/go-eagle/eagle/cmd/eagle/internal/cache"
	"github.com/go-eagle/eagle/cmd/eagle/internal/project"
	"github.com/go-eagle/eagle/cmd/eagle/internal/run"
	"github.com/go-eagle/eagle/cmd/eagle/internal/upgrade"

	"github.com/spf13/cobra"
)

var (
	// Version is the version of the compiled software.
	Version = "v0.5.0"

	rootCmd = &cobra.Command{
		Use:     "eagle",
		Short:   "Eagle: An develop kit for Go microservices.",
		Long:    `Eagle: An develop kit for Go microservices.`,
		Version: Version,
	}
)

func init() {
	rootCmd.AddCommand(project.CmdNew)
	rootCmd.AddCommand(run.CmdRun)
	rootCmd.AddCommand(cache.CmdCache)
	rootCmd.AddCommand(proto.CmdProto)
	rootCmd.AddCommand(upgrade.CmdUpgrade)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
