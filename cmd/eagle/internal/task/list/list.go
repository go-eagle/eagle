package list

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"

	"github.com/go-eagle/eagle/cmd/eagle/internal/utils"
)

// CmdAdd represents the new command.
var CmdList = &cobra.Command{
	Use:   "list",
	Short: "List all task",
	Long:  "List all task. Example: eagle task list",
	Run:   run,
}

var (
	targetDir string
)

func init() {
	CmdList.Flags().StringVarP(&targetDir, "-target-dir", "t", "internal/tasks", "generate target directory")
}

func run(cmd *cobra.Command, args []string) {
	if len(targetDir) == 0 {
		fmt.Println(color.RedString("Please enter the target dir"))
		return
	}
	// eg: eagle task list
	f, err := os.Open(targetDir)
	if err != nil {
		fmt.Printf("%s\n", color.RedString(err.Error()))
		return
	}
	files, err := f.ReadDir(-1)
	_ = f.Close()
	if err != nil {
		return
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Task Name", "Handler Name", "File Name", "Location"})

	for key, file := range files {
		fileName := file.Name()
		// 下划线转冒号分隔
		taskName := strings.Replace(fileName, "_", ":", 5)
		taskName = strings.Split(taskName, ".")[0]
		handlerName := fmt.Sprintf("Handle%sTask", strings.Split(utils.Case2Camel(fileName), ".")[0])
		t.AppendRow([]interface{}{key + 1, taskName, handlerName, fileName, targetDir})
	}
	t.Render()
}
