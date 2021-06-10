package upgrade

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/1024casts/snake/cmd/snake/internal/base"
)

// CmdUpgrade represents the upgrade command.
var CmdUpgrade = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade the snake tools",
	Long:  "Upgrade the snake tools. Example: snake upgrade",
	Run:   Run,
}

// Run upgrade the snake tools.
func Run(cmd *cobra.Command, args []string) {
	err := base.GoGet(
		"github.com/1024casts/snake/cmd/snake",
	)
	if err != nil {
		fmt.Println(err)
	}
}
