package log

import (
	"testing"

	"github.com/1024casts/snake/config"
)

func TestInitLogger(t *testing.T) {
	// init config
	if err := config.Init("../../conf/config.local.yaml"); err != nil {
		panic(err)
	}

	logger := InitLogger()

	logger.Info("test logger info")
	logger.Warn("test logger warning")
	logger.Error("test logger error")
}
