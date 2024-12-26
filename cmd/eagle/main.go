package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/go-eagle/eagle/cmd/eagle/internal/cache"
	"github.com/go-eagle/eagle/cmd/eagle/internal/handler"
	"github.com/go-eagle/eagle/cmd/eagle/internal/model"
	"github.com/go-eagle/eagle/cmd/eagle/internal/project"
	"github.com/go-eagle/eagle/cmd/eagle/internal/proto"
	"github.com/go-eagle/eagle/cmd/eagle/internal/repo"
	"github.com/go-eagle/eagle/cmd/eagle/internal/run"
	"github.com/go-eagle/eagle/cmd/eagle/internal/service"
	"github.com/go-eagle/eagle/cmd/eagle/internal/task"
	"github.com/go-eagle/eagle/cmd/eagle/internal/upgrade"
)

var (
	// Version is the version of the compiled software.
	Version = "v1.0.3"

	rootCmd = &cobra.Command{
		Use:     "eagle",
		Short:   "Eagle: A microservice framework for Go",
		Long:    `Eagle: A microservice framework for Go`,
		Version: Version,
	}
)

func init() {
	rootCmd.AddCommand(project.CmdNew)
	rootCmd.AddCommand(run.CmdRun)
	rootCmd.AddCommand(handler.CmdHandler)
	rootCmd.AddCommand(cache.CmdCache)
	rootCmd.AddCommand(repo.CmdRepo)
	rootCmd.AddCommand(service.CmdService)
	rootCmd.AddCommand(proto.CmdProto)
	rootCmd.AddCommand(task.CmdTask)
	rootCmd.AddCommand(model.CmdNew)
	rootCmd.AddCommand(upgrade.CmdUpgrade)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
