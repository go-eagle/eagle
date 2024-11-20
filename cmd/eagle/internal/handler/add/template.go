package add

import (
	"bytes"
	"strings"

	"github.com/alecthomas/template"
)

const handlerTemplate = `
package v1

import (
    "github.com/gin-gonic/gin"
    "github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/errcode"

	"github.com/go-eagle/eagle-layout/internal/service"
	"github.com/go-eagle/eagle-layout/internal/types"
)

// {{.Name}} {{.LcName}}
// @Summary {{.LcName}}
// @Description {{.LcName}}
// @Tags system
// @Accept  json
// @Produce  json
// @Router /{{.UsName}} {{.Method}}
func {{.Name}}(c *gin.Context) {
	var req types.{{.Name}}Request
	{{- if .Method eq "GET" }}
	if err := c.ShouldBindQuery(&req); err != nil {
		app.Error(c, errcode.ErrInvalidParam.WithDetails(err.Error()))
		return
	}
	{{- end }}

	{{- if .Method eq "POST" }}
	if err := c.ShouldBindJSON(&req); err != nil {
		app.Error(c, errcode.ErrInvalidParam.WithDetails(err.Error()))
		return
	}
	{{- end }}

	{{- if .Method eq "PUT" }}
	if err := c.ShouldBindJSON(&req); err != nil {
		app.Error(c, errcode.ErrInvalidParam.WithDetails(err.Error()))
		return
	}
	{{- end }}

	{{- if .Method eq "PATCH" }}
	if err := c.ShouldBindJSON(&req); err != nil {
		app.Error(c, errcode.ErrInvalidParam.WithDetails(err.Error()))
		return
	}
	{{- end }}

	{{- if .Method eq "DELETE" }}
	if err := c.ShouldBindJSON(&req); err != nil {
		app.Error(c, errcode.ErrInvalidParam.WithDetails(err.Error()))
		return
	}
	{{- end }}

	var ret any
	// change to your service
	// ret, err := service.GreeterSvc.Hello(c, req.Name)
	// if err != nil {
	// 	app.Error(c, err)
	// 	return
	// }

	app.Success(c, ret)
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
