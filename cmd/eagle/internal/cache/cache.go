package cache

import (
	"github.com/go-eagle/eagle/cmd/eagle/internal/cache/add"
	"github.com/spf13/cobra"
)

// CmdProto represents the proto command.
var CmdCache = &cobra.Command{
	Use:   "cache",
	Short: "Generate the cache file",
	Long:  "Generate the cache file.",
	Run:   run,
}

func init() {
	CmdCache.AddCommand(add.CmdAdd)
}

func run(cmd *cobra.Command, args []string) {
}
