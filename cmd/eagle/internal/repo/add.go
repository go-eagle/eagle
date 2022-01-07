package repo

import (
	"fmt"

	"github.com/go-eagle/eagle/cmd/eagle/internal/utils"

	"github.com/spf13/cobra"
)

// CmdRepo represents the new command.
var CmdRepo = &cobra.Command{
	Use:   "repo",
	Short: "Create a repo by template",
	Long:  "Create a repo using the repo template. Example: eagle repo UserCache",
	Run:   run,
}

var (
	targetDir string
)

func init() {
	CmdRepo.Flags().StringVarP(&targetDir, "-target-dir", "t", "internal/repository", "generate target directory")
}

func run(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Please enter the cache filename")
		return
	}
	// eg: eagle cache UserCache
	filename := args[0]

	c := &Repo{
		Name:    utils.Ucfirst(filename), // 首字母大写
		Path:    targetDir,
		ModName: utils.ModName(),
	}
	if err := c.Generate(); err != nil {
		fmt.Println(err)
		return
	}
}
