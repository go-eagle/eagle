package add

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/mod/modfile"
)

// CmdAdd represents the add command.
var CmdAdd = &cobra.Command{
	Use:   "add",
	Short: "Add a proto API template",
	Long:  "Add a proto API template. eg: eagle add hello.proto",
	Run:   run,
}

// eg: eagle proto add helloworld/v1/helloworld.proto
func run(cmd *cobra.Command, args []string) {
	input := args[0]
	n := strings.LastIndex(input, "/")
	if n == -1 {
		fmt.Println("The proto path needs to be hierarchical.")
		return
	}
	path := input[:n]
	fileName := input[n+1:]
	pkgName := strings.ReplaceAll(path, "/", ".")

	p := &Proto{
		Name:        fileName,
		Path:        path,
		Package:     pkgName,
		GoPackage:   goPackage(path),
		JavaPackage: javaPackage(pkgName),
		Service:     serviceName(fileName),
	}
	if err := p.Generate(); err != nil {
		fmt.Println(err)
		return
	}
}

func modName() string {
	modBytes, err := os.ReadFile("go.mod")
	if err != nil {
		if modBytes, err = os.ReadFile("../go.mod"); err != nil {
			return ""
		}
	}
	return modfile.ModulePath(modBytes)
}

func goPackage(path string) string {
	s := strings.Split(path, "/")
	return modName() + "/" + path + ";" + s[len(s)-1]
}

func javaPackage(name string) string {
	return name
}

func serviceName(name string) string {
	return export(strings.Split(name, ".")[0])
}

func export(s string) string { return strings.ToUpper(s[:1]) + s[1:] }
