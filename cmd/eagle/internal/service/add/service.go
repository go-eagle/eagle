package add

import (
	"fmt"
	"os"
	"path"

	"github.com/go-eagle/eagle/cmd/eagle/internal/utils"
)

// Service is a svc generator.
type Service struct {
	Name    string
	LcName  string
	Path    string
	Service string
	Package string
	ModName string
}

// Generate generate a svc template.
func (s *Service) Generate() error {
	body, err := s.execute()
	if err != nil {
		return err
	}
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	to := path.Join(wd, s.Path)
	if _, err := os.Stat(to); os.IsNotExist(err) {
		if err := os.MkdirAll(to, 0700); err != nil {
			return err
		}
	}
	name := path.Join(to, utils.Camel2Case(s.Name)+"_svc.go")
	if _, err := os.Stat(name); !os.IsNotExist(err) {
		return fmt.Errorf("%s already exists", s.Name)
	}
	return os.WriteFile(name, body, 0644)
}
