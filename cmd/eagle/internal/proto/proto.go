package proto

import (
	"github.com/go-eagle/eagle/cmd/eagle/internal/proto/add"

	"github.com/spf13/cobra"
)

// CmdProto represents the proto command.
var CmdProto = &cobra.Command{
	Use:   "proto",
	Short: "Generate the proto files",
	Long:  "Generate the proto files.",
	Run:   run,
}

func init() {
	CmdProto.AddCommand(add.CmdAdd)
}

func run(cmd *cobra.Command, args []string) {
}
