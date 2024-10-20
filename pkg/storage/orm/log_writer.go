package orm

import (
	"fmt"

	"github.com/go-eagle/eagle/pkg/log"
	"gorm.io/gorm/logger"
)

type LoggerWriter struct {
	log log.Logger
}

func NewLogWriter(log log.Logger) logger.Writer {
	return &LoggerWriter{
		log: log,
	}
}

func (l *LoggerWriter) Printf(s string, v ...interface{}) {
	l.log.Info(fmt.Sprintf(s, v...))
}
