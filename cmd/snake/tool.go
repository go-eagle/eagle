package main

import (
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	//"github.com/fatih/color"
)

// Tool is snake tool.
type Tool struct {
	Name         string    `json:"name"`
	Alias        string    `json:"alias"`
	BuildTime    time.Time `json:"build_time"`
	Install      string    `json:"install"`
	Requirements []string  `json:"requirements"`
	Dir          string    `json:"dir"`
	Summary      string    `json:"summary"`
	Platform     []string  `json:"platform"`
	Author       string    `json:"author"`
	URL          string    `json:"url"`
	Hidden       bool      `json:"hidden"`
	requires     []*Tool
}

func (t Tool) needUpdated() bool {
	for _, r := range t.requires {
		if r.needUpdated() {
			return true
		}
	}
	if !t.supportOS() || t.Install == "" {
		return false
	}
	if f, err := os.Stat(t.toolPath()); err == nil {
		if t.BuildTime.After(f.ModTime()) {
			return true
		}
	}
	return false
}

func (t Tool) toolPath() string {
	name := t.Alias
	if name == "" {
		name = t.Name
	}
	gobin := Getenv("GOBIN")
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	if gobin != "" {
		return filepath.Join(gobin, name)
	}
	return filepath.Join(gopath(), "bin", name)
}

func (t Tool) installed() bool {
	_, err := os.Stat(t.toolPath())
	return err == nil
}

func (t Tool) supportOS() bool {
	for _, p := range t.Platform {
		if strings.ToLower(p) == runtime.GOOS {
			return true
		}
	}
	return false
}

func (t Tool) install() {
	if t.Install == "" {
		//fmt.Fprintf(os.Stderr, color.RedString("%s: 自动安装失败详情请查看文档：%s\n", t.Name, toolDoc))
		return
	}
	fmt.Println(t.Install)
	cmds := strings.Split(t.Install, " ")
	if len(cmds) > 0 {
		if err := runTool(t.Name, path.Dir(t.toolPath()), cmds[0], cmds[1:]); err == nil {
			//color.Green("%s: 安装成功!", t.Name)
		}
	}
}

func (t Tool) updated() bool {
	if !t.supportOS() || t.Install == "" {
		return false
	}
	if f, err := os.Stat(t.toolPath()); err == nil {
		if t.BuildTime.After(f.ModTime()) {
			return true
		}
	}
	return false
}

func toolList() (tools []*Tool) {
	return toolIndexs
}

func installAndRun(name string, args []string) (err error) {
	for _, t := range toolList() {
		if name == t.Name {
			if !t.installed() || t.needUpdated() {
				t.install()
			}
			pwd, _ := os.Getwd()
			err = runTool(t.Name, pwd, t.toolPath(), args)
			return
		}
	}
	return fmt.Errorf("找不到%s", name)
}

func gopath() (gp string) {
	gopaths := strings.Split(Getenv("GOPATH"), string(filepath.ListSeparator))

	if len(gopaths) == 1 && gopaths[0] != "" {
		return gopaths[0]
	}
	pwd, err := os.Getwd()
	if err != nil {
		return
	}
	abspwd, err := filepath.Abs(pwd)
	if err != nil {
		return
	}
	for _, gopath := range gopaths {
		if gopath == "" {
			continue
		}
		absgp, err := filepath.Abs(gopath)
		if err != nil {
			return
		}
		if strings.HasPrefix(abspwd, absgp) {
			return absgp
		}
	}
	return build.Default.GOPATH
}

func runTool(name, dir, cmd string, args []string) (err error) {
	toolCmd := &exec.Cmd{
		Path:   cmd,
		Args:   append([]string{cmd}, args...),
		Dir:    dir,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Env:    os.Environ(),
	}
	if filepath.Base(cmd) == cmd {
		var lp string
		if lp, err = exec.LookPath(cmd); err == nil {
			toolCmd.Path = lp
		}
	}
	if err = toolCmd.Run(); err != nil {
		if e, ok := err.(*exec.ExitError); !ok || !e.Exited() {
			fmt.Fprintf(os.Stderr, "运行 %s 出错: %v\n", name, err)
		}
	}
	return
}
