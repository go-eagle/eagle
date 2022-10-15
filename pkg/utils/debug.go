package utils

import (
	"runtime"

	"github.com/go-eagle/eagle/pkg/log"
)

// PrintStackTrace print stack info
func PrintStackTrace(msg string, err interface{}) {
	buf := make([]byte, 64*1024)
	buf = buf[:runtime.Stack(buf, false)]
	log.Error("%s, err: %s\nstack: %s", msg, err, buf)
}
