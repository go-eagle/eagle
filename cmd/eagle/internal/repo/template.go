package repo

import (
	"bytes"
	"strings"

	"github.com/alecthomas/template"
)

const repoTemplate = `

`

func (r *Repo) execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("repo").Parse(strings.TrimSpace(repoTemplate))
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, r); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
