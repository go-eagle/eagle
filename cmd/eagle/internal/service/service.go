package service

import (
	"github.com/spf13/cobra"

	"github.com/go-eagle/eagle/cmd/eagle/internal/service/add"
)

// CmdService represents the service command.
var CmdService = &cobra.Command{
	Use:     "svc",
	Aliases: []string{"service"},
	Short:   "Generate the service/svc file",
	Long:    "Generate the service/svc file.",
	Run:     run,
}

func init() {
	CmdService.AddCommand(add.CmdAdd)
}

func run(cmd *cobra.Command, args []string) {
}
