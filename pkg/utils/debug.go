package utils

import (
	"fmt"
	"runtime"
)

// PrintStackTrace print stack info
func PrintStackTrace(msg string, err interface{}) string {
	buf := make([]byte, 64*1024)
	buf = buf[:runtime.Stack(buf, false)]
	return fmt.Sprintf("%s, err: %s\nstack: %s", msg, err, buf)
}
