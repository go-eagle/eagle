package add

import (
	"fmt"
	"os"
	"path"

	"github.com/go-eagle/eagle/cmd/eagle/internal/utils"
)

// Handler is a cache generator.
type Handler struct {
	Name    string
	LcName  string
	UsName  string
	Path    string
	Version string
	Method  string
}

// Generate generate a Handler template.
func (h *Handler) Generate() error {
	body, err := h.execute()
	if err != nil {
		return err
	}
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	to := path.Join(wd, h.Path, h.Version)
	if _, err := os.Stat(to); os.IsNotExist(err) {
		if err := os.MkdirAll(to, 0700); err != nil {
			return err
		}
	}
	name := path.Join(to, utils.Camel2Case(h.Name)+".go")
	if _, err := os.Stat(name); !os.IsNotExist(err) {
		return fmt.Errorf("%s already exists", h.Name)
	}
	return os.WriteFile(name, body, 0644)
}
