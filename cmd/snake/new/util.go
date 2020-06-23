package new

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gobuffalo/packr/v2"
	"github.com/urfave/cli"

	"github.com/1024casts/snake/cmd/snake/common"
	"github.com/1024casts/snake/pkg/util/color"
)

// CreateProject create a template project for Jupiter
func CreateProject(cli *cli.Context) (err error) {
	newArgs := cli.Args()
	if len(newArgs) <= 0 {
		fmt.Println(color.Red("Command line new execution error, please use snake new -h for details"))
		return
	}
	name := newArgs[0]
	if name == "" {
		project.Name = DefaultProjectName
	} else {
		project.Name = name
	}
	if project.Path != "" {
		if project.Path, err = filepath.Abs(project.Path); err != nil {
			return
		}
		project.Path = filepath.Join(project.Path, project.Name)
	} else {
		pwd, _ := os.Getwd()
		project.Path = filepath.Join(pwd, project.Name)
	}
	//GetModPath
	modPath := common.GetModPath(project.Path)
	fmt.Println("new project modPrefix:", modPath)
	project.ModPrefix = modPath
	if err = doCreateProject(); err != nil {
		return
	}
	fmt.Println(color.Greenf("Project dir:", project.Path))
	fmt.Println(color.Green("Project created successfully"))
	return
}

//go:generate packr2
func doCreateProject() (err error) {
	box := packr.New("all", "./templates")
	if err = os.MkdirAll(project.Path, 0755); err != nil {
		return
	}
	for _, name := range box.List() {
		if project.ModPrefix != "" && name == "go.mod.tmpl" {
			continue
		}
		tmpl, _ := box.FindString(name)
		i := strings.LastIndex(name, string(os.PathSeparator))
		if i > 0 {
			dir := name[:i]
			if err = os.MkdirAll(filepath.Join(project.Path, dir), 0755); err != nil {
				return
			}
		}
		if strings.HasSuffix(name, ".tmpl") {
			name = strings.TrimSuffix(name, ".tmpl")
		}
		if err = doWriteFile(filepath.Join(project.Path, name), tmpl); err != nil {
			return
		}
	}

	return
}

func doWriteFile(path, tmpl string) (err error) {
	data, err := parseTmpl(tmpl)
	if err != nil {
		return
	}
	fmt.Println(color.Greenf("File generated----------------------->", path))
	return ioutil.WriteFile(path, data, 0755)
}

func parseTmpl(tmpl string) ([]byte, error) {
	tmp, err := template.New("").Parse(tmpl)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err = tmp.Execute(&buf, project); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
