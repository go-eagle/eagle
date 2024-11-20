package add

import (
	"bytes"
	"html/template"
	"strings"
)

const repoTemplate = `
package repository

//go:generate mockgen -source={{.UsName}}_repo.go -destination=../../internal/mocks/{{.UsName}}_repo_mock.go  -package mocks

import (
	"context"
	"fmt"
	"strings"
	"time"

	localCache "github.com/go-eagle/eagle/pkg/cache"
	"github.com/go-eagle/eagle/pkg/encoding"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"

	"{{.ModName}}/internal/dal"
	"{{.ModName}}/internal/dal/cache"
	"{{.ModName}}/internal/dal/db/dao"
	"{{.ModName}}/internal/dal/db/model"
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
	db         *dal.DBClient
	tracer 		 trace.Tracer
{{- if eq .WithCache true }}
	cache  		 cache.{{.Name}}Cache
{{- end }}
	localCache localCache.Cache
	sg         singleflight.Group
}

{{- if .WithCache }}
// New{{.Name}} new a repository and return
func New{{.Name}}(db *dal.DBClient, cache cache.{{.Name}}Cache) {{.Name}}Repo {
	return &{{.LcName}}Repo{
		db:     		db,
		tracer: 		otel.Tracer("{{.LcName}}"),
		cache:  		cache,
		localCache: localCache.NewMemoryCache("local:{{.LcName}}:", encoding.JSONEncoding{}),
		sg:         singleflight.Group{},
	}
}
{{- else }}
// New{{.Name}} new a repository and return
func New{{.Name}}(db *dal.DBClient) {{.Name}}Repo {
	return &{{.LcName}}Repo{
		db:     db,
		tracer: otel.Tracer("{{.LcName}}Repo"),
	}
}
{{- end }}

// Create{{.Name}} create a item
func (r *{{.LcName}}Repo) Create{{.Name}}(ctx context.Context, data *model.{{.Name}}Model) (id int64, err error) {
	err = dao.{{.Name}}Model.WithContext(ctx).Create(&data).Error
	if err != nil {
		return 0, errors.Wrap(err, "[repo] create {{.Name}} err")
	}

	return data.ID, nil
}

// Update{{.Name}} update item
func (r *{{.LcName}}Repo) Update{{.Name}}(ctx context.Context, id int64, data *model.{{.Name}}Model) error {
	_, err := dao.{{.Name}}Model.WithContext(ctx).Where(dao.{{.Name}}Model.ID.Eq(id)).Updates(data)
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
	// read db
	data := new(model.{{.Name}}Model)
	err = r.db.WithContext(ctx).Raw(fmt.Sprintf(_get{{.Name}}SQL, _table{{.Name}}Name), id).Scan(&data).Error
	if err != nil {
		return
	}

{{- if .WithCache }}
	// write cache
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
	// get data from local cache
	err = r.localCache.Get(ctx, cast.ToString(id), &ret)
	if err != nil {
		return nil, err
	}
	if ret != nil && ret.ID > 0 {
		return ret, nil
	}

	// read redis cache
	ret, err = r.cache.Get{{.Name}}Cache(ctx, id)
	if err != nil {
		return nil, err
	}
	if ret != nil && ret.ID > 0 {
		return ret, nil
	}

	// get data from db
	// 避免缓存击穿(瞬间有大量请求过来)
	val, err, _ := r.sg.Do("sg:{{.LcName}}:"+cast.ToString(id), func() (interface{}, error) {
		// read db or rpc
		data, err := dao.{{.Name}}Model.WithContext(ctx).Where(dao.{{.Name}}Model.ID.Eq(id)).First()
		if err != nil {
			// cache not found and set empty cache to avoid 缓存穿透
			// Note: 如果缓存空数据太多，会大大降低缓存命中率，可以改为使用布隆过滤器
			if errors.Is(err, gorm.ErrRecordNotFound) {
				r.cache.SetCacheWithNotFound(ctx, id)
			}
			return nil, errors.Wrapf(err, "[repo] get {{.Name}} from db error, id: %d", id)
		}

		// write cache
		if data != nil && data.ID > 0 {
			// write redis
			err = r.cache.Set{{.Name}}Cache(ctx, id, data, 5*time.Minute)
			if err != nil {
				return nil, errors.WithMessage(err, "[repo] Get{{.Name}} Set{{.Name}}Cache error")
			}

			// write local cache
			err = r.localCache.Set(ctx, cast.ToString(id), data, 2*time.Minute)
			if err != nil {
				return nil, errors.WithMessage(err, "[repo] Get{{.Name}} localCache set error")
			}
		}
		return data, nil
	})
	if err != nil {
		return nil, err
	}

	return val.(*model.{{.Name}}Model), nil
{{- else }}
	// read db
	ret, err = dao.{{.Name}}Model.WithContext(ctx).Where(dao.{{.Name}}Model.ID.In(ids...)).Find()
	if err != nil {
		return nil, err
	}
	return ret, nil
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
