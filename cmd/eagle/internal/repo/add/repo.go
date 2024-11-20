package add

import (
	"fmt"
	"os"
	"path"

	"github.com/go-eagle/eagle/cmd/eagle/internal/utils"
)

// Repo is a cache generator.
type Repo struct {
	Name      string
	LcName    string
	UsName    string
	Path      string
	Service   string
	Package   string
	ModName   string
	WithCache bool
}

// Generate generate a repo template.
func (r *Repo) Generate() error {
	body, err := r.execute()
	if err != nil {
		return err
	}
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	to := path.Join(wd, r.Path)
	if _, err := os.Stat(to); os.IsNotExist(err) {
		if err := os.MkdirAll(to, 0700); err != nil {
			return err
		}
	}
	name := path.Join(to, utils.Camel2Case(r.Name)+"_repo.go")
	if _, err := os.Stat(name); !os.IsNotExist(err) {
		return fmt.Errorf("%s already exists", r.Name)
	}
	return os.WriteFile(name, body, 0644)
}
