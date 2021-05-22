package sql

import (
	"github.com/1024casts/snake/pkg/breaker"
	// database driver
	_ "github.com/go-sql-driver/mysql"

	xtime "github.com/1024casts/snake/pkg/time"
)

// Config mysql config.
type Config struct {
	DSN          string          // write data source name.
	ReadDSN      []string        // read data source name.
	Active       int             // pool
	Idle         int             // pool
	IdleTimeout  xtime.Duration  // connect max life time.
	QueryTimeout xtime.Duration  // query sql timeout
	ExecTimeout  xtime.Duration  // execute sql timeout
	TranTimeout  xtime.Duration  // transaction sql timeout
	Breaker      *breaker.Config // breaker
}

// NewMySQL new db and retry connection when has error.
//func NewMySQL(c *Config) (db *DB) {
//	if c.QueryTimeout == 0 || c.ExecTimeout == 0 || c.TranTimeout == 0 {
//		panic("mysql must be set query/execute/transction timeout")
//	}
//	db, err := Open(c)
//	if err != nil {
//		log.Error("open mysql error(%v)", err)
//		panic(err)
//	}
//	return
//}
