package sql

import (
	"context"
	"database/sql"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

// Stmt prepared stmt.
type Stmt struct {
	db    *conn
	tx    bool
	query string
	stmt  atomic.Value
	trace trace.Tracer
}

// Close closes the statement.
func (s *Stmt) Close() (err error) {
	if s == nil {
		err = ErrStmtNil
		return
	}
	stmt, ok := s.stmt.Load().(*sql.Stmt)
	if ok {
		err = errors.WithStack(stmt.Close())
	}
	return
}

// Exec executes a prepared statement with the given arguments and returns a
// Result summarizing the effect of the statement.
func (s *Stmt) Exec(ctx context.Context, args ...interface{}) (res sql.Result, err error) {
	if s == nil {
		err = ErrStmtNil
		return
	}
	now := time.Now()
	defer slowLog(fmt.Sprintf("Exec query(%s) args(%+v)", s.query, args), now)
	tr := otel.Tracer("sql")
	if s.trace != nil {
		tr = s.trace
	}
	ctx, span := tr.Start(ctx, "Exec")
	defer span.End()

	span.SetAttributes(
		semconv.DBSystemMySQL,
		attribute.String("db.instance", s.db.addr),
		semconv.PeerServiceKey.String("mysql"),
		semconv.DBStatementKey.String(fmt.Sprintf("exec: %s, args: %+v", s.query, args)),
	)
	if err = s.db.breaker.Allow(); err != nil {
		_metricReqErr.Inc(s.db.addr, s.db.addr, "stmt:exec", "breaker")
		return
	}
	stmt, ok := s.stmt.Load().(*sql.Stmt)
	if !ok {
		err = ErrStmtNil
		return
	}
	_, c, cancel := s.db.conf.ExecTimeout.Shrink(ctx)
	res, err = stmt.ExecContext(c, args...)
	cancel()
	s.db.onBreaker(&err)
	_metricReqDur.Observe(int64(time.Since(now)/time.Millisecond), s.db.addr, s.db.addr, "stmt:exec")
	if err != nil {
		err = errors.Wrapf(err, "exec: %s, args: %+v", s.query, args)
	}
	return
}

// Query executes a prepared query statement with the given arguments and
// returns the query results as a *Rows.
func (s *Stmt) Query(ctx context.Context, args ...interface{}) (rows *Rows, err error) {
	if s == nil {
		err = ErrStmtNil
		return
	}
	now := time.Now()
	defer slowLog(fmt.Sprintf("Query query: %s, args: %+v", s.query, args), now)
	tr := otel.Tracer("sql")
	if s.trace != nil {
		tr = s.trace
	}
	ctx, span := tr.Start(ctx, "Query")
	defer span.End()

	span.SetAttributes(
		semconv.DBSystemMySQL,
		attribute.String("db.instance", s.db.addr),
		semconv.PeerServiceKey.String("mysql"),
		semconv.DBStatementKey.String(fmt.Sprintf("exec: %s, args: %+v", s.query, args)),
	)
	if err = s.db.breaker.Allow(); err != nil {
		_metricReqErr.Inc(s.db.addr, s.db.addr, "stmt:query", "breaker")
		return
	}
	stmt, ok := s.stmt.Load().(*sql.Stmt)
	if !ok {
		err = ErrStmtNil
		return
	}
	_, c, cancel := s.db.conf.QueryTimeout.Shrink(ctx)
	// nolint: rowserrcheck
	rs, err := stmt.QueryContext(c, args...)
	s.db.onBreaker(&err)
	_metricReqDur.Observe(int64(time.Since(now)/time.Millisecond), s.db.addr, s.db.addr, "stmt:query")
	if err != nil {
		err = errors.Wrapf(err, "query: %s, args: %+v", s.query, args)
		cancel()
		return
	}
	rows = &Rows{Rows: rs, cancel: cancel}
	return
}

// QueryRow executes a prepared query statement with the given arguments.
// If an error occurs during the execution of the statement, that error will
// be returned by a call to Scan on the returned *Row, which is always non-nil.
// If the query selects no rows, the *Row's Scan will return ErrNoRows.
// Otherwise, the *Row's Scan scans the first selected row and discards the rest.
func (s *Stmt) QueryRow(ctx context.Context, args ...interface{}) (row *Row) {
	now := time.Now()
	defer slowLog(fmt.Sprintf("QueryRow query: %s, args: %+v", s.query, args), now)
	row = &Row{db: s.db, query: s.query, args: args}
	if s == nil {
		row.err = ErrStmtNil
		return
	}
	tr := otel.Tracer("sql")
	if s.trace != nil {
		tr = s.trace
	}
	ctx, span := tr.Start(ctx, "QueryRow")
	defer span.End()

	span.SetAttributes(
		semconv.DBSystemMySQL,
		attribute.String("db.instance", s.db.addr),
		semconv.PeerServiceKey.String("mysql"),
		semconv.DBStatementKey.String(fmt.Sprintf("exec: %s, args: %+v", s.query, args)),
	)
	if row.err = s.db.breaker.Allow(); row.err != nil {
		_metricReqErr.Inc(s.db.addr, s.db.addr, "stmt:queryRow", "breaker")
		return
	}
	stmt, ok := s.stmt.Load().(*sql.Stmt)
	if !ok {
		return
	}
	_, c, cancel := s.db.conf.QueryTimeout.Shrink(ctx)
	row.Row = stmt.QueryRowContext(c, args...)
	row.cancel = cancel
	_metricReqDur.Observe(int64(time.Since(now)/time.Millisecond), s.db.addr, s.db.addr, "stmt:queryRow")
	return
}
