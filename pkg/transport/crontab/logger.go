package crontab

import (
	"strings"
	"time"

	"github.com/go-eagle/eagle/pkg/log"
)

type Logger struct {
	Log log.Logger
}

func (l Logger) Info(msg string, keysAndValues ...interface{}) {
	keysAndValues = formatTimes(keysAndValues)
	keysAndValues = append([]interface{}{
		msg,
	}, keysAndValues...)
	l.Log.Infof(formatString(len(keysAndValues)), keysAndValues...)
}

func (l Logger) Error(err error, msg string, keysAndValues ...interface{}) {
	keysAndValues = formatTimes(keysAndValues)
	keysAndValues = append([]interface{}{
		msg,
		"error", err,
	}, keysAndValues...)
	l.Log.Errorf(formatString(len(keysAndValues)+2), keysAndValues...)
}

// formatString returns a logfmt-like format string for the number of
// key/values.
func formatString(numKeysAndValues int) string {
	var sb strings.Builder
	sb.WriteString("%s")
	if numKeysAndValues > 0 {
		sb.WriteString(", ")
	}
	for i := 0; i < numKeysAndValues/2; i++ {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString("%v=%v")
	}
	return sb.String()
}

// formatTimes formats any time.Time values as RFC3339.
func formatTimes(keysAndValues []interface{}) []interface{} {
	var formattedArgs []interface{}
	for _, arg := range keysAndValues {
		if t, ok := arg.(time.Time); ok {
			arg = t.Format(time.RFC3339)
		}
		formattedArgs = append(formattedArgs, arg)
	}
	return formattedArgs
}
