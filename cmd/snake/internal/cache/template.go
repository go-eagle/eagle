package cache

import (
	"bytes"
	"strings"

	"github.com/alecthomas/template"
)

const cacheTemplate = `

`

type Cache struct {
}

func (p *Cache) execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("cache").Parse(strings.TrimSpace(cacheTemplate))
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, p); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
