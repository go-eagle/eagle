package cache

import (
	"fmt"

	"github.com/go-eagle/eagle/cmd/eagle/internal/utils"

	"github.com/spf13/cobra"
)

// CmdCache represents the new command.
var CmdCache = &cobra.Command{
	Use:   "cache",
	Short: "Create a cache by template",
	Long:  "Create a cache using the cache template. Example: eagle cache UserCache",
	Run:   run,
}

var (
	targetDir string
)

func init() {
	CmdCache.Flags().StringVarP(&targetDir, "-target-dir", "t", "internal/cache", "generate target directory")
}

func run(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Please enter the cache filename")
		return
	}
	// eg: eagle cache UserCache
	filename := args[0]

	c := &Cache{
		Name:    utils.Ucfirst(filename), // 首字母大写
		Path:    targetDir,
		ModName: utils.ModName(),
	}
	if err := c.Generate(); err != nil {
		fmt.Println(err)
		return
	}
}
