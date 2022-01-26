package add

import (
	"bytes"
	"strings"

	"github.com/alecthomas/template"
)

var svcTemplate = `
package service

import (
	"context"

	"{{.ModName}}/internal/repository"
)

// {{.Name}}Service define a interface
type {{.Name}}Service interface {
	Hello(ctx context.Context) error
}

type {{.LcName}}Service struct {
	repo repository.Repository
}

var _ {{.Name}}Service = (*{{.LcName}}Service)(nil)

func new{{.Name}}Service(svc *service) *{{.LcName}}Service {
	return &{{.LcName}}Service{repo: svc.repo}
}

// Hello .
func (s *{{.LcName}}Service) Hello(ctx context.Context) error {
	return nil
}
`

func (s *Service) execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("service").Parse(strings.TrimSpace(svcTemplate))
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
