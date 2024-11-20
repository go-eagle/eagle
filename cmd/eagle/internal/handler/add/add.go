package add

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/go-eagle/eagle/cmd/eagle/internal/utils"
)

// CmdCache represents the new command.
var CmdAdd = &cobra.Command{
	Use:   "add",
	Short: "Create a handler file by template",
	Long:  "Create a handler file using the handler template. Example: eagle handler add demo",
	Run:   run,
}

var (
	targetDir string
	version   string
	method    string
)

func init() {
	CmdAdd.Flags().StringVarP(&targetDir, "target-dir", "t", "internal/handler", "generate target directory")
	CmdAdd.Flags().StringVarP(&version, "version", "v", "v1", "handler version")
	CmdAdd.Flags().StringVarP(&method, "method", "m", "GET", "http method, eg: GET, POST, PUT, PATCH, DELETE")
}

func run(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Please enter the handler filename")
		return
	}
	// eg: eagle handler add demo
	filename := args[0]

	c := &Handler{
		Name:    utils.Ucfirst(filename),    // 首字母大写
		LcName:  utils.Lcfirst(filename),    // 首字母小写
		UsName:  utils.Camel2Case(filename), // 下划线分隔
		Path:    targetDir,
		Version: version,
		Method:  method,
	}
	if err := c.Generate(); err != nil {
		fmt.Println(err)
		return
	}
}
