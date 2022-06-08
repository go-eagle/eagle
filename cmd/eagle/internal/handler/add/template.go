package add

import (
	"bytes"
	"strings"

	"github.com/alecthomas/template"
)

const handlerTemplate = `
package handler

import (
    "github.com/gin-gonic/gin"
    "github.com/go-eagle/eagle/pkg/app"
)

// {{.Name}} {{.LcName}}
// @Summary {{.LcName}}
// @Description {{.LcName}}
// @Tags system
// @Accept  json
// @Produce  json
// @Router /{{.UsName}} {{.Method}}
func {{.Name}}(c *gin.Context) {
    // here add your code

    app.Success(c, gin.H{})
}
`

func (h *Handler) execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("handler").Parse(strings.TrimSpace(handlerTemplate))
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, h); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
