package add

import (
	"bytes"
	"html/template"
	"strings"
)

const repoTemplate = `
package repository

//go:generate mockgen -source=internal/repository/{{.UsName}}_repo.go -destination=internal/mock/{{.UsName}}_repo_mock.go  -package mock

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"

	"{{.ModName}}/internal/cache"
	"{{.ModName}}/internal/model"
)

var (
	_table{{.Name}}Name   = (&model.{{.Name}}Model{}).TableName()
	_get{{.Name}}SQL      = "SELECT * FROM %s WHERE id = ?"
	_batchGet{{.Name}}SQL = "SELECT * FROM %s WHERE id IN (%s)"
)

var _ {{.Name}}Repo = (*{{.LcName}}Repo)(nil)

// {{.Name}}Repo define a repo interface
type {{.Name}}Repo interface {
	Create{{.Name}}(ctx context.Context, data *model.{{.Name}}Model) (id int64, err error)
	Update{{.Name}}(ctx context.Context, id int64, data *model.{{.Name}}Model) error
	Get{{.Name}}(ctx context.Context, id int64) (ret *model.{{.Name}}Model, err error)
	BatchGet{{.Name}}(ctx context.Context, ids []int64) (ret []*model.{{.Name}}Model, err error)
}

type {{.LcName}}Repo struct {
	db     *gorm.DB
	tracer trace.Tracer
{{- if eq .WithCache true }}
	cache  *cache.{{.Name}}Cache
{{- end }}
}

{{- if .WithCache }}
// New{{.Name}} new a repository and return
func New{{.Name}}(db *gorm.DB, cache *cache.{{.Name}}Cache) {{.Name}}Repo {
	return &userBaseRepo{
		db:     db,
		tracer: otel.Tracer("userBaseRepo"),
		cache:  cache,
	}
}
{{- else }}
// New{{.Name}} new a repository and return
func New{{.Name}}(db *gorm.DB) {{.Name}}Repo {
	return &{{.LcName}}Repo{
		db:     db,
		tracer: otel.Tracer("{{.LcName}}Repo"),
	}
}
{{- end }}

// Create{{.Name}} create a item
func (r *{{.LcName}}Repo) Create{{.Name}}(ctx context.Context, data *model.{{.Name}}Model) (id int64, err error) {
	err = r.db.WithContext(ctx).Create(&data).Error
	if err != nil {
		return 0, errors.Wrap(err, "[repo] create {{.Name}} err")
	}

	return data.ID, nil
}

// Update{{.Name}} update item
func (r *{{.LcName}}Repo) Update{{.Name}}(ctx context.Context, id int64, data *model.{{.Name}}Model) error {
	item, err := r.Get{{.Name}}(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "[repo] update {{.Name}} err: %v", err)
	}
	err = r.db.Model(&item).Updates(data).Error
	if err != nil {
		return err
	}

{{- if .WithCache }}
	// delete cache
	_ = r.cache.Del{{.Name}}Cache(ctx, id)
{{- end }}
	return nil
}

// Get{{.Name}} get a record
func (r *{{.LcName}}Repo) Get{{.Name}}(ctx context.Context, id int64) (ret *model.{{.Name}}Model, err error) {
{{- if .WithCache }}
	// read cache
	item, err := r.cache.Get{{.Name}}Cache(ctx, id)
	if err != nil {
		return nil, err
	}
	if item != nil {
		return item, nil
	}
{{- end }}
	data := new(model.{{.Name}}Model)
	err = r.db.WithContext(ctx).Raw(fmt.Sprintf(_get{{.Name}}SQL, _table{{.Name}}Name), id).Scan(&data).Error
	if err != nil {
		return
	}
s
{{- if .WithCache }}
	if data.ID > 0 {
		err = r.cache.Set{{.Name}}Cache(ctx, id, data, 5*time.Minute)
		if err != nil {
			return nil, err
		}
	}
{{- end }}
	return data, nil
}

// BatchGet{{.Name}} batch get items
func (r *{{.LcName}}Repo) BatchGet{{.Name}}(ctx context.Context, ids []int64) (ret []*model.{{.Name}}Model, err error) {
{{- if .WithCache }}
	idsStr := cast.ToStringSlice(ids)
	itemMap, err := r.cache.MultiGet{{.Name}}Cache(ctx, ids)
	if err != nil {
		return nil, err
	}
	var missedID []int64
	for _, v := range ids {
		item, ok := itemMap[cast.ToString(v)]
		if !ok {
			missedID = append(missedID, v)
			continue
		}
		ret = append(ret, item)
	}
	// get missed data
	if len(missedID) > 0 {
		var missedData []*model.{{.Name}}Model
		_sql := fmt.Sprintf(_batchGet{{.Name}}SQL, _table{{.Name}}Name, strings.Join(idsStr, ","))
		err = r.db.WithContext(ctx).Raw(_sql).Scan(&missedData).Error
		if err != nil {
			// you can degrade to ignore error
			return nil, err
		}
		if len(missedData) > 0 {
			ret = append(ret, missedData...)
			err = r.cache.MultiSet{{.Name}}Cache(ctx, missedData, 5*time.Minute)
			if err != nil {
				// you can degrade to ignore error
				return nil, err
			}
		}
	}
	return ret, nil
{{- else }}
	items := make([]*model.{{.Name}}Model, 0)
	err = r.db.WithContext(ctx).Raw(fmt.Sprintf(_batchGet{{.Name}}SQL, _table{{.Name}}Name), ids).Scan(&items).Error
	if err != nil {
		return
	}
	return items, nil
{{- end }}
}
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
