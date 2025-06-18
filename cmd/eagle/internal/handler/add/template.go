package add

import (
	"bytes"
	"html/template"
	"strings"
)

const handlerTemplate = `
package v1

import (
    "github.com/gin-gonic/gin"
    "github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/errcode"

	// "{{.ModName}}/internal/service"
	"{{.ModName}}/internal/types"
)

// {{.Name}}Handler {{.LcName}}
type {{.Name}}Handler struct {
	// here you can add your service
	// example:
	// 	UserService service.UserService
}

// New{{.Name}}Handler create a new {{.Name}}Handler
func New{{.Name}}Handler() *{{.Name}}Handler {
	return &{{.Name}}Handler{}
}

// {{.Name}} {{.LcName}}
// @Summary {{.LcName}}
// @Description {{.LcName}}
// @Tags system
// @Accept  json
// @Produce  json
// @Router /{{.UsName}} {{.Method}}
func (h *{{.Name}}Handler) {{.Name}}(c *gin.Context) {
	var req types.{{.Name}}Request
	{{- if eq .Method "GET" }}
	if err := c.ShouldBindQuery(&req); err != nil {
		app.Error(c, errcode.ErrInvalidParam.WithDetails(err.Error()))
		return
	}
	{{- end }}

	{{- if eq .Method "POST" }}
	if err := c.ShouldBindJSON(&req); err != nil {
		app.Error(c, errcode.ErrInvalidParam.WithDetails(err.Error()))
		return
	}
	{{- end }}

	{{- if eq .Method "PUT" }}
	if err := c.ShouldBindJSON(&req); err != nil {
		app.Error(c, errcode.ErrInvalidParam.WithDetails(err.Error()))
		return
	}
	{{- end }}

	{{- if eq .Method "PATCH" }}
	if err := c.ShouldBindJSON(&req); err != nil {
		app.Error(c, errcode.ErrInvalidParam.WithDetails(err.Error()))
		return
	}
	{{- end }}

	{{- if eq .Method "DELETE" }}
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
