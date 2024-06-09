package model

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/go-eagle/eagle/cmd/eagle/internal/base"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
)

// Model is a model template.
type Model struct {
	Filename    string
	Host        string
	MysqlPort   string
	Database    string
	TableName   string
	PackageName string
	StructName  string
	User        string
	Password    string
	Format      string
	TargetDir   string
}

// New new a project from remote repo.
func (m *Model) Generate() error {
	to := path.Join(m.TargetDir, m.Filename)
	if _, err := os.Stat(to); !os.IsNotExist(err) {
		fmt.Printf("🚫 %s already exists\n", m.Filename)
		override := false
		prompt := &survey.Confirm{
			Message: "📂 Do you want to override the file ?",
			Help:    "Delete the existing file and create the file.",
		}
		e := survey.AskOne(prompt, &override)
		if e != nil {
			return e
		}
		if !override {
			return err
		}
		e = os.RemoveAll(to)
		if e != nil {
			return e
		}
	}

	var (
		err error
	)
	// check db2struct if installed
	if err = look("db2struct"); err != nil {
		// install db2struct
		err := base.GoInstall(
			"github.com/Shelnutt2/db2struct/cmd/db2struct",
		)
		if err != nil {
			return err
		}
	}

	fmt.Printf("🚀 Creating model %s, please wait a moment.\n", m.Filename)

	input := []string{
		"--host", m.Host,
		"--mysql_port", m.MysqlPort,
		"--database", m.Database,
		"--table", m.TableName,
		"--package", m.PackageName,
		"--struct", m.StructName,
		"--user", m.User,
		"-p", m.Password,
		"--gorm",
	}

	// default to add json
	if len(m.Format) == 0 {
		input = append(input, "--json")
	} else {
		// --no-json
		input = append(input, m.Format)
	}

	// read stdout to buf
	var (
		stdout bytes.Buffer
	)
	cmd := exec.Command("db2struct", input...)
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		_, err = fmt.Fprintf(os.Stderr, "\033[31mERROR: failed to generate model(%s)\033[m\n", err.Error())
		return err
	}

	// write stdout to file
	if err = os.WriteFile(to, stdout.Bytes(), 0644); err != nil {
		_, err = fmt.Fprintf(os.Stderr, "\033[31mERROR: write to model(%s) error\033[m\n", err.Error())
		return err
	}

	fmt.Printf("\n🍺 Model creation succeeded %s\n", color.GreenString(m.Filename))
	return nil
}

// look find executable file
func look(name ...string) error {
	for _, n := range name {
		if _, err := exec.LookPath(n); err != nil {
			return err
		}
	}
	return nil
}
