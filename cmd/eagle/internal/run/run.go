package run

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

// CmdRun run project command.
var CmdRun = &cobra.Command{
	Use:   "run",
	Short: "Run project",
	Long:  "Run project. Example: eagle run",
	Run:   Run,
}

// Run project.
func Run(cmd *cobra.Command, args []string) {
	var dir string
	cmdArgs, programArgs := splitArgs(cmd, args)
	if len(cmdArgs) > 0 {
		dir = cmdArgs[0]
	}
	base, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err)
		return
	}
	if dir == "" {
		// find the directory containing the cmd/*
		cmdPath, err := findCMD(base)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err)
			return
		}
		if len(cmdPath) == 0 {
			fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", "The cmd directory cannot be found in the current directory")
			return
		} else if len(cmdPath) == 1 {
			for _, v := range cmdPath {
				dir = v
			}
		} else {
			var cmdPaths []string
			for k, _ := range cmdPath {
				cmdPaths = append(cmdPaths, k)
			}
			prompt := &survey.Select{
				Message:  "Which directory do you want to run?",
				Options:  cmdPaths,
				PageSize: 6,
			}
			// get cmd dir, eg: cmd/server
			err := survey.AskOne(prompt, &dir)
			if err != nil && dir == "" {
				return
			}

			// project absolute path
			dir = cmdPath[dir]
		}
	}

	// go run /absolute/path/cmd/server
	fd := exec.Command("go", append([]string{"run", dir}, programArgs...)...)
	fd.Env = os.Environ()
	fd.Stdout = os.Stdout
	fd.Stderr = os.Stderr
	fd.Dir = dir
	if err := fd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err.Error())
		return
	}
}

// ./cmd subcmd cmdarg1 cmdarg2 -- progamarg3 progamarg4
// example:
// | Command                                                       | ArgsLenAtDash |
// |---------------------------------------------------------------|---------------|
// | `./cmd subcmd arg1 arg2 -- arg3 arg4`                         | 2             |
// | `./cmd subcmd -- arg1 arg2 arg3 arg4`                         | 0             |
// | `./cmd subcmd arg1 arg2 arg3 arg4`                            | -1            |
// | `./cmd subcmd arg1 --flag=f arg2 --otherflag=o -- arg3 arg4 --not-flag=12` | 2             |
func splitArgs(cmd *cobra.Command, args []string) (cmdArgs, programArgs []string) {
	dashAt := cmd.ArgsLenAtDash()
	// have string --, or dashAt == -1
	if dashAt >= 0 {
		return args[:dashAt], args[dashAt:]
	}
	// only return cmd args
	return args, []string{}
}

// map, eg: cmd/server -> project absolute path
func findCMD(base string) (map[string]string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	if !strings.HasSuffix(wd, "/") {
		wd += "/"
	}
	var root bool
	// return map[cmd/server:/eagle-example/cmd/server cmd/consumer:/eagle-example/cmd/consumer
	next := func(dir string) (map[string]string, error) {
		cmdPath := make(map[string]string)
		err := filepath.Walk(dir, func(walkPath string, info os.FileInfo, err error) error {
			// multi level directory is not allowed under the cmdPath directory, so it is judged that the path ends with cmdPath.
			if strings.HasSuffix(walkPath, "cmd") {
				paths, err := os.ReadDir(walkPath)
				if err != nil {
					return err
				}
				for _, fileInfo := range paths {
					if fileInfo.IsDir() {
						abs := filepath.Join(walkPath, fileInfo.Name())
						cmdPath[strings.TrimPrefix(abs, wd)] = abs
					}
				}
				return nil
			}
			if info.Name() == "go.mod" {
				root = true
			}
			return nil
		})
		return cmdPath, err
	}

	for i := 0; i < 5; i++ {
		tmp := base
		cmd, err := next(tmp)
		if err != nil {
			return nil, err
		}
		if len(cmd) > 0 {
			return cmd, nil
		}
		if root {
			break
		}
		_ = filepath.Join(base, "..")
	}
	return map[string]string{"": base}, nil
}
