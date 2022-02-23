package add

import (
	"fmt"

	"github.com/go-eagle/eagle/cmd/eagle/internal/utils"

	"github.com/spf13/cobra"
)

// CmdAdd represents the new command.
var CmdAdd = &cobra.Command{
	Use:   "add",
	Short: "Create a repo file by template",
	Long:  "Create a repo file using the repo template. Example: eagle repo add UserCache",
	Run:   run,
}

var (
	targetDir string
	withCache bool
)

func init() {
	CmdAdd.Flags().StringVarP(&targetDir, "-target-dir", "t", "internal/repository", "generate target directory")
	CmdAdd.Flags().BoolVarP(&withCache, "-with-cache", "c", true, "add cache operate")
}

func run(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Please enter the repo filename")
		return
	}
	// eg: eagle cache UserCache
	filename := args[0]

	c := &Repo{
		Name:      utils.Ucfirst(filename),    // 首字母大写
		LcName:    utils.Lcfirst(filename),    // 首字母小写
		UsName:    utils.Camel2Case(filename), // 下划线分隔
		Path:      targetDir,
		ModName:   utils.ModName(),
		WithCache: withCache,
	}
	if err := c.Generate(); err != nil {
		fmt.Println(err)
		return
	}
}
