package upgrade

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/go-eagle/eagle/cmd/eagle/internal/base"
)

// CmdUpgrade represents the upgrade command.
var CmdUpgrade = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade the eagle tools",
	Long:  "Upgrade the eagle tools. Example: eagle upgrade",
	Run:   Run,
}

// Run upgrade the eagle tools.
func Run(cmd *cobra.Command, args []string) {
	err := base.GoInstall(
		"github.com/go-eagle/eagle/cmd/eagle",
		"github.com/go-eagle/eagle/cmd/protoc-gen-go-gin",
		"google.golang.org/protobuf/cmd/protoc-gen-go",
		"google.golang.org/grpc/cmd/protoc-gen-go-grpc",
		"github.com/envoyproxy/protoc-gen-validate",
		"github.com/google/gnostic",
		"github.com/google/gnostic/cmd/protoc-gen-openapi",
		"github.com/google/wire/cmd/wire",
	)
	if err != nil {
		fmt.Println(err)
	}
}
