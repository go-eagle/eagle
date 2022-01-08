package add

import (
	"bytes"
	"strings"

	"github.com/alecthomas/template"
)

const repoTemplate = `
package repository

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"{{.ModName}}/internal/model"
)

var (
	_tableName           = (&model.{{.Name}}Model{}).TableName()
	_get{{.Name}}SQL      = "SELECT * FROM %s WHERE id = ?"
	_batchGet{{.Name}}SQL = "SELECT * FROM %s WHERE id IN (?)"
)

// Create{{.Name}} create a item
func (r *repository) Create{{.Name}}(ctx context.Context, data *model.{{.Name}}Model) (id int64, err error) {
	err = r.db.WithContext(ctx).Create(&data).Error
	if err != nil {
		return 0, errors.Wrap(err, "[repo] create {{.Name}} err")
	}

	return data.ID, nil
}

// Update{{.Name}} update item
func (r *repository) Update{{.Name}}(ctx context.Context, id int64, data *model.{{.Name}}Model) error {
	item, err := r.Get{{.Name}}(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "[repo] update {{.Name}} err: %v", err)
	}
	err = r.db.Model(&item).Updates(data).Error
	if err != nil {
		return err
	}
	return nil
}

// Get{{.Name}} get a record
func (r *repository) Get{{.Name}}(ctx context.Context, id int64) (ret *model.{{.Name}}Model, err error) {
	item := new(model.{{.Name}}Model)
	err = r.db.WithContext(ctx).Raw(fmt.Sprintf(_get{{.Name}}SQL, _tableName), id).Scan(&item).Error
	if err != nil {
		return
	}
	return item, nil
}

// BatchGet{{.Name}} batch get items
func (r *repository) BatchGet{{.Name}}(ctx context.Context, ids int64) (ret []*model.{{.Name}}Model, err error) {
	items := make([]*model.{{.Name}}Model, 0)
	err = r.db.WithContext(ctx).Raw(fmt.Sprintf(_batchGet{{.Name}}SQL, _tableName), ids).Scan(&items).Error
	if err != nil {
		return
	}
	return items, nil
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
