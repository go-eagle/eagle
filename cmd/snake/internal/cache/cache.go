package cache

import (
	"os"

	"github.com/spf13/cobra"
)

// CmdCache represents the new command.
var CmdCache = &cobra.Command{
	Use:   "cache",
	Short: "Create a cache by template",
	Long:  "Create a cache using the cache template. Example: snake cache -name UserCache",
	Run:   run,
}

var repoURL string

func init() {
	if repoURL = os.Getenv("SNAKE_LAYOUT_REPO"); repoURL == "" {
		repoURL = "https://github.com/1024casts/snake-layout.git"
	}
	CmdCache.Flags().StringVarP(&repoURL, "-repo-url", "r", repoURL, "layout repo")
}

func run(cmd *cobra.Command, args []string) {

}
