package handler

import (
	"github.com/go-eagle/eagle/cmd/eagle/internal/handler/add"
	"github.com/spf13/cobra"
)

// CmdHandler represents the handler command.
var CmdHandler = &cobra.Command{
	Use:   "handler",
	Short: "Generate the handler file",
	Long:  "Generate the handler/controller file.",
	Run:   run,
}

func init() {
	CmdHandler.AddCommand(add.CmdAdd)
}

func run(cmd *cobra.Command, args []string) {
}
