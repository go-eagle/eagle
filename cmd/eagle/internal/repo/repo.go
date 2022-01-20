package repo

import (
	"github.com/spf13/cobra"

	"github.com/go-eagle/eagle/cmd/eagle/internal/repo/add"
)

// CmdProto represents the proto command.
var CmdRepo = &cobra.Command{
	Use:   "repo",
	Short: "Generate the repo file",
	Long:  "Generate the repo file.",
	Run:   run,
}

func init() {
	CmdRepo.AddCommand(add.CmdAdd)
}

func run(cmd *cobra.Command, args []string) {
}
