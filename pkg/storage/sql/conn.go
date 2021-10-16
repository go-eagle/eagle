package sql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	breaker "github.com/go-kratos/aegis/circuitbreaker"
)

// conn database connection
type conn struct {
	*sql.DB
	breaker breaker.CircuitBreaker
	conf    *Config
	addr    string
}

func (db *conn) onBreaker(err *error) {
	if err != nil && *err != nil && *err != sql.ErrNoRows && *err != sql.ErrTxDone {
		db.breaker.MarkFailed()
	} else {
		db.breaker.MarkSuccess()
	}
}

func (db *conn) begin(ctx context.Context) (tx *Tx, err error) {
	now := time.Now()
	defer slowLog("Begin", now)
	tr := otel.Tracer("sql")
	ctx, span := tr.Start(ctx, "conn.begin")
	defer span.End()
	span.SetAttributes(
		semconv.DBSystemMySQL,
		attribute.String("db.instance", db.addr),
	)

	if err = db.breaker.Allow(); err != nil {
		_metricReqErr.Inc(db.addr, db.addr, "begin", err.Error())
		return
	}
	_, c, cancel := db.conf.TranTimeout.Shrink(ctx)
	rtx, err := db.BeginTx(c, nil)
	_metricReqDur.Observe(int64(time.Since(now)/time.Millisecond), db.addr, db.addr, "begin")
	if err != nil {
		err = errors.WithStack(err)
		cancel()
		return
	}
	tx = &Tx{tx: rtx, trace: tr, db: db, c: c, cancel: cancel}
	return
}

func (db *conn) exec(ctx context.Context, query string, args ...interface{}) (res sql.Result, err error) {
	now := time.Now()
	defer slowLog(fmt.Sprintf("Exec query: %s args: %+v", query, args), now)
	ctx, span := otel.Tracer("sql").Start(ctx, "conn.exec")
	defer span.End()

	span.SetAttributes(
		semconv.DBSystemMySQL,
		attribute.String("db.instance", db.addr),
		semconv.DBStatementKey.String(fmt.Sprintf("exec: %s, args: %+v", query, args)),
	)

	if err = db.breaker.Allow(); err != nil {
		_metricReqErr.Inc(db.addr, db.addr, "exec", "breaker")
		return
	}
	_, c, cancel := db.conf.ExecTimeout.Shrink(ctx)
	res, err = db.ExecContext(c, query, args...)
	cancel()
	db.onBreaker(&err)
	_metricReqDur.Observe(int64(time.Since(now)/time.Millisecond), db.addr, db.addr, "exec")
	if err != nil {
		err = errors.Wrapf(err, "exec: %s, args: %+v", query, args)
	}
	return
}

func (db *conn) ping(ctx context.Context) (err error) {
	now := time.Now()
	defer slowLog("Ping", now)
	ctx, span := otel.Tracer("sql").Start(ctx, "conn.ping")
	defer span.End()

	span.SetAttributes(
		semconv.DBSystemMySQL,
		attribute.String("db.instance", db.addr),
	)

	if err = db.breaker.Allow(); err != nil {
		_metricReqErr.Inc(db.addr, db.addr, "ping", "breaker")
		return
	}
	_, c, cancel := db.conf.ExecTimeout.Shrink(ctx)
	err = db.PingContext(c)
	cancel()
	db.onBreaker(&err)
	_metricReqDur.Observe(int64(time.Since(now)/time.Millisecond), db.addr, db.addr, "ping")
	if err != nil {
		err = errors.WithStack(err)
	}
	return
}

func (db *conn) prepare(query string) (*Stmt, error) {
	defer slowLog(fmt.Sprintf("Prepare query(%s)", query), time.Now())
	stmt, err := db.Prepare(query)
	if err != nil {
		err = errors.Wrapf(err, "prepare %s", query)
		return nil, err
	}
	st := &Stmt{query: query, db: db}
	st.stmt.Store(stmt)
	return st, nil
}

func (db *conn) prepared(query string) (stmt *Stmt) {
	defer slowLog(fmt.Sprintf("Prepared query(%s)", query), time.Now())
	stmt = &Stmt{query: query, db: db}
	s, err := db.Prepare(query)
	if err == nil {
		stmt.stmt.Store(s)
		return
	}
	go func() {
		for {
			s, err := db.Prepare(query)
			if err != nil {
				time.Sleep(time.Second)
				continue
			}
			stmt.stmt.Store(s)
			return
		}
	}()
	return
}

func (db *conn) query(ctx context.Context, query string, args ...interface{}) (rows *Rows, err error) {
	now := time.Now()
	defer slowLog(fmt.Sprintf("Query query: %s args: %+v", query, args), now)
	ctx, span := otel.Tracer("sql").Start(ctx, "query")
	defer span.End()

	span.SetAttributes(
		semconv.DBSystemMySQL,
		attribute.String("db.instance", db.addr),
		semconv.PeerServiceKey.String("mysql"),
		semconv.DBStatementKey.String(fmt.Sprintf("exec: %s, args: %+v", query, args)),
	)

	if err = db.breaker.Allow(); err != nil {
		_metricReqErr.Inc(db.addr, db.addr, "query", "breaker")
		return
	}
	_, c, cancel := db.conf.QueryTimeout.Shrink(ctx)
	// nolint: rowserrcheck
	rs, err := db.DB.QueryContext(c, query, args...)
	db.onBreaker(&err)
	_metricReqDur.Observe(int64(time.Since(now)/time.Millisecond), db.addr, db.addr, "query")
	if err != nil {
		err = errors.Wrapf(err, "query: %s, args: %+v", query, args)
		cancel()
		return
	}
	rows = &Rows{Rows: rs, cancel: cancel}
	return
}

func (db *conn) queryRow(ctx context.Context, query string, args ...interface{}) *Row {
	now := time.Now()
	defer slowLog(fmt.Sprintf("QueryRow query: %s args: %+v", query, args), now)
	tr := otel.Tracer("sql")
	ctx, span := tr.Start(ctx, "queryRow")
	defer span.End()

	span.SetAttributes(
		semconv.DBSystemMySQL,
		attribute.String("db.instance", db.addr),
		semconv.PeerServiceKey.String("mysql"),
		semconv.DBStatementKey.String(fmt.Sprintf("exec: %s, args: %+v", query, args)),
	)

	if err := db.breaker.Allow(); err != nil {
		_metricReqErr.Inc(db.addr, db.addr, "queryRow", "breaker")
		return &Row{db: db, trace: nil, err: err}
	}
	_, c, cancel := db.conf.QueryTimeout.Shrink(ctx)
	r := db.DB.QueryRowContext(c, query, args...)
	_metricReqDur.Observe(int64(time.Since(now)/time.Millisecond), db.addr, db.addr, "queryRow")
	return &Row{db: db, Row: r, query: query, args: args, trace: tr, cancel: cancel}
}
