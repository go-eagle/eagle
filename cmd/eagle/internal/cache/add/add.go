package add

import (
	"fmt"

	"github.com/go-eagle/eagle/cmd/eagle/internal/utils"

	"github.com/spf13/cobra"
)

// CmdCache represents the new command.
var CmdAdd = &cobra.Command{
	Use:   "add",
	Short: "Create a cache file by template",
	Long:  "Create a cache file using the cache template. Example: eagle cache add UserCache",
	Run:   run,
}

var (
	targetDir string
)

func init() {
	CmdAdd.Flags().StringVarP(&targetDir, "-target-dir", "t", "internal/cache", "generate target directory")
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
