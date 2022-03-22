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
	"log"

	"github.com/hibiken/asynq"
)

const (
	// Type{{.Name}} task name
	Type{{.Name}} = "{{.ColonName}}"
)

// {{.Name}}Payload define data payload
type {{.Name}}Payload struct {
	UserID int
}

// New{{.Name}}Task to create a task. 
func New{{.Name}}Task(userID int) (*asynq.Task, error) {
	payload, err := json.Marshal({{.Name}}Payload{UserID: userID})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(Type{{.Name}}, payload), nil
}

// Handle{{.Name}}Task to handle the input task.
func Handle{{.Name}}Task(ctx context.Context, t *asynq.Task) error {
	var p {{.Name}}Payload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
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
