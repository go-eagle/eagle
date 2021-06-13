package cache

import (
	"fmt"
	"unicode"

	"github.com/spf13/cobra"
)

// CmdCache represents the new command.
var CmdCache = &cobra.Command{
	Use:   "cache",
	Short: "Create a cache by template",
	Long:  "Create a cache using the cache template. Example: snake cache -name UserCache",
	Run:   run,
}

var (
	targetDir string
)

func init() {
	CmdCache.Flags().StringVarP(&targetDir, "-target-dir", "t", "internal/cache", "generate target directory")
}

func run(cmd *cobra.Command, args []string) {
	// eg: snake cache UserLike
	filename := args[0]

	c := &Cache{
		Name: Ucfirst(filename), // 首字母大写
		Path: targetDir,
	}
	if err := c.Generate(); err != nil {
		fmt.Println(err)
		return
	}
}

// 首字母大写
func Ucfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}
