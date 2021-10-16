package sql

import (
	"database/sql"
	"time"

	breaker "github.com/go-kratos/aegis/circuitbreaker"
	"github.com/go-kratos/aegis/circuitbreaker/sre"
	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"

	// database driver
	_ "github.com/go-sql-driver/mysql"

	"github.com/go-eagle/eagle/pkg/container/group"
	"github.com/go-eagle/eagle/pkg/log"
	xtime "github.com/go-eagle/eagle/pkg/time"
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

// Config mysql config.
type Config struct {
	DSN             string         // write data source name.
	ReadDSN         []string       // read data source name.
	MaxOpenConn     int            // open pool
	MaxIdleConn     int            // idle pool
	ConnMaxLifeTime xtime.Duration // connect max life time.
	QueryTimeout    xtime.Duration // query sql timeout
	ExecTimeout     xtime.Duration // execute sql timeout
	TranTimeout     xtime.Duration // transaction sql timeout
}

// NewMySQL new db and retry connection when has error.
func NewMySQL(c *Config) (db *DB) {
	if c.QueryTimeout == 0 || c.ExecTimeout == 0 || c.TranTimeout == 0 {
		panic("mysql must be set query/execute/transction timeout")
	}
	db, err := Open(c)
	if err != nil {
		log.Error("open mysql error(%v)", err)
		panic(err)
	}
	return
}

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
	brkGroup := group.NewGroup(func() interface{} {
		return sre.NewBreaker()
	})
	brk := brkGroup.Get(addr)
	w := &conn{DB: d, breaker: brk.(breaker.CircuitBreaker), conf: c, addr: addr}
	rs := make([]*conn, 0, len(c.ReadDSN))
	for _, rd := range c.ReadDSN {
		d, err := connect(c, rd)
		if err != nil {
			return nil, err
		}
		addr = parseDSNAddr(rd)
		brk := brkGroup.Get(addr)
		r := &conn{DB: d, breaker: brk.(breaker.CircuitBreaker), conf: c, addr: addr}
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
	d.SetMaxOpenConns(c.MaxOpenConn)
	d.SetMaxIdleConns(c.MaxIdleConn)
	d.SetConnMaxLifetime(time.Duration(c.ConnMaxLifeTime))
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
