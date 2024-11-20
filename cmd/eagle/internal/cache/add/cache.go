package add

import (
	"fmt"
	"os"
	"path"

	"github.com/go-eagle/eagle/cmd/eagle/internal/utils"
)

// Cache is a cache generator.
type Cache struct {
	Name      string
	LcName    string
	UsName    string
	ColonName string
	Path      string
	Service   string
	Package   string
	ModName   string
}

// Generate generate a cache template.
func (c *Cache) Generate() error {
	body, err := c.execute()
	if err != nil {
		return err
	}
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	to := path.Join(wd, c.Path)
	if _, err := os.Stat(to); os.IsNotExist(err) {
		if err := os.MkdirAll(to, 0700); err != nil {
			return err
		}
	}
	name := path.Join(to, utils.Camel2Case(c.Name)+"_cache.go")
	if _, err := os.Stat(name); !os.IsNotExist(err) {
		return fmt.Errorf("%s already exists", c.Name)
	}
	return os.WriteFile(name, body, 0644)
}
