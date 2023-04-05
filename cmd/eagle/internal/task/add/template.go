package add

import (
	"bytes"
	"html/template"
	"strings"
)

const taskTemplate = `
package tasks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"

	"github.com/go-eagle/eagle/pkg/log"
	"github.com/pkg/errors"
)

const (
	// Type{{.Name}} task name
	Type{{.Name}} = "{{.ColonName}}"
)

// {{.Name}}Payload define data payload
type {{.Name}}Payload struct {
	// fill your biz field
}

// New{{.Name}}Task to create a task. 
func New{{.Name}}Task(data {{.Name}}Payload) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return errors.Wrapf(err, "[tasks] json marshal error, name: %s", Type{{.Name}})
	}
	task := asynq.NewTask(Type{{.Name}}, payload)
	_, err = GetClient().Enqueue(task)
	if err != nil {
		return errors.Wrapf(err, "[tasks] Enqueue task error, name: %s", Type{{.Name}})
	}

	return nil
}

// Handle{{.Name}}Task to handle the input task.
func Handle{{.Name}}Task(ctx context.Context, t *asynq.Task) error {
	var p {{.Name}}Payload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		log.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	// here to write biz logic

	return nil
}
`

func (t *Task) execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("task").Parse(strings.TrimSpace(taskTemplate))
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, t); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
