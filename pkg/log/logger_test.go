package log

import (
	"testing"

	"github.com/1024casts/snake/config"
)

func TestWriteLog(t *testing.T) {
	// init config
	if err := config.Init("../../conf/config.local.yaml"); err != nil {
		panic(err)
	}

	l := InitLogger()
	l.Warn("test warn")
	l.Error("test error")
	l.Info("test info")
}
