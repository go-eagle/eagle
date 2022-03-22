package add

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/go-eagle/eagle/cmd/eagle/internal/utils"
)

// Repo is a cache generator.
type Task struct {
	Name      string
	LcName    string
	UsName    string
	ColonName string
	Path      string
	Service   string
	Package   string
	ModName   string
	WithCache bool
}

// Generate generate a repo template.
func (t *Task) Generate() error {
	body, err := t.execute()
	if err != nil {
		return err
	}
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	to := path.Join(wd, t.Path)
	if _, err := os.Stat(to); os.IsNotExist(err) {
		if err := os.MkdirAll(to, 0700); err != nil {
			return err
		}
	}
	name := path.Join(to, utils.Camel2Case(t.Name)+".go")
	if _, err := os.Stat(name); !os.IsNotExist(err) {
		return fmt.Errorf("%s already exists", t.Name)
	}
	return ioutil.WriteFile(name, body, 0644)
}
