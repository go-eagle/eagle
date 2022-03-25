package task

import (
	"github.com/go-eagle/eagle/cmd/eagle/internal/task/list"
	"github.com/spf13/cobra"

	"github.com/go-eagle/eagle/cmd/eagle/internal/task/add"
)

// CmdProto represents the proto command.
var CmdTask = &cobra.Command{
	Use:   "task",
	Short: "Generate the task file",
	Long:  "Generate the task file.",
	Run:   run,
}

func init() {
	CmdTask.AddCommand(add.CmdAdd)
	CmdTask.AddCommand(list.CmdList)
}

func run(cmd *cobra.Command, args []string) {
}
