package sql

import (
	"database/sql"
	"time"

	"github.com/1024casts/snake/pkg/breaker"
	"github.com/1024casts/snake/pkg/log"
	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

const (
	_family          = "sql_client"
	_slowLogDuration = time.Millisecond * 250
)

var (
	// ErrStmtNil prepared stmt error
	ErrStmtNil = errors.New("sql: prepare failed and stmt nil")
	// ErrNoMaster is returned by Master when call master multiple times.
	ErrNoMaster = errors.New("sql: no master instance")
	// ErrNoRows is returned by Scan when QueryRow doesn'trace return a row.
	// In such a case, QueryRow returns a placeholder *Row value that defers
	// this error until a Scan.
	ErrNoRows = sql.ErrNoRows
	// ErrTxDone transaction done.
	ErrTxDone = sql.ErrTxDone
)

// Open opens a database specified by its database driver name and a
// driver-specific data source name, usually consisting of at least a database
// name and connection information.
func Open(c *Config) (*DB, error) {
	db := new(DB)
	d, err := connect(c, c.DSN)
	if err != nil {
		return nil, err
	}
	addr := parseDSNAddr(c.DSN)
	brkGroup := breaker.NewGroup(c.Breaker)
	brk := brkGroup.Get(addr)
	w := &conn{DB: d, breaker: brk, conf: c, addr: addr}
	rs := make([]*conn, 0, len(c.ReadDSN))
	for _, rd := range c.ReadDSN {
		d, err := connect(c, rd)
		if err != nil {
			return nil, err
		}
		addr = parseDSNAddr(rd)
		brk := brkGroup.Get(addr)
		r := &conn{DB: d, breaker: brk, conf: c, addr: addr}
		rs = append(rs, r)
	}
	db.write = w
	db.read = rs
	db.master = &DB{write: db.write}
	return db, nil
}

func connect(c *Config, dataSourceName string) (*sql.DB, error) {
	d, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		err = errors.WithStack(err)
		return nil, err
	}
	d.SetMaxOpenConns(c.Active)
	d.SetMaxIdleConns(c.Idle)
	d.SetConnMaxLifetime(time.Duration(c.IdleTimeout))
	return d, nil
}

// parseDSNAddr parse dsn name and return addr.
func parseDSNAddr(dsn string) (addr string) {
	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		// just ignore parseDSN error, mysql client will return error for us when connect.
		return ""
	}
	return cfg.Addr
}

func slowLog(statement string, now time.Time) {
	du := time.Since(now)
	if du > _slowLogDuration {
		log.Warn("%s slow log statement: %s time: %v", _family, statement, du)
	}
}
