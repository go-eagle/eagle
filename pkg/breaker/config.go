package breaker

import (
	"time"

	xtime "github.com/go-eagle/eagle/pkg/time"
)

// Config broker config.
type Config struct {
	SwitchOff bool // breaker switch,default off.

	// Google
	K float64

	Window  xtime.Duration
	Bucket  int
	Request int64
}

func (conf *Config) fix() {
	if conf.K == 0 {
		conf.K = 1.5
	}
	if conf.Request == 0 {
		conf.Request = 100
	}
	if conf.Bucket == 0 {
		conf.Bucket = 10
	}
	if conf.Window == 0 {
		conf.Window = xtime.Duration(3 * time.Second)
	}
}
