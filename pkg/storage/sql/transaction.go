package sql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

// Tx transaction.
type Tx struct {
	db     *conn
	tx     *sql.Tx
	trace  trace.Tracer
	c      context.Context
	cancel func()
}

// Commit commits the transaction.
func (tx *Tx) Commit() (err error) {
	err = tx.tx.Commit()
	tx.cancel()
	tx.db.onBreaker(&err)
	if tx.trace != nil {
		_, span := tx.trace.Start(tx.c, "tx.Commit")
		defer span.End()

		span.SetAttributes(
			semconv.DBSystemMySQL,
			attribute.String("db.instance", tx.db.addr),
			semconv.PeerServiceKey.String("mysql"),
		)
	}
	if err != nil {
		err = errors.WithStack(err)
	}
	return
}

// Rollback aborts the transaction.
func (tx *Tx) Rollback() (err error) {
	err = tx.tx.Rollback()
	tx.cancel()
	tx.db.onBreaker(&err)
	if tx.trace != nil {
		_, span := tx.trace.Start(tx.c, "tx.Rollback")
		defer span.End()

		span.SetAttributes(
			semconv.DBSystemMySQL,
			attribute.String("db.instance", tx.db.addr),
			semconv.PeerServiceKey.String("mysql"),
		)
	}
	if err != nil {
		err = errors.WithStack(err)
	}
	return
}

// Exec executes a query that doesn'trace return rows. For example: an INSERT and
// UPDATE.
func (tx *Tx) Exec(query string, args ...interface{}) (res sql.Result, err error) {
	now := time.Now()
	defer slowLog(fmt.Sprintf("Exec query: %s, args: %+v", query, args), now)

	if tx.trace != nil {
		_, span := tx.trace.Start(tx.c, "tx.Exec")
		defer span.End()

		span.SetAttributes(
			semconv.DBSystemMySQL,
			attribute.String("db.instance", tx.db.addr),
			semconv.PeerServiceKey.String("mysql"),
			semconv.DBStatementKey.String(fmt.Sprintf("exec: %s, args: %+v", query, args)),
		)
	}
	res, err = tx.tx.ExecContext(tx.c, query, args...)
	_metricReqDur.Observe(int64(time.Since(now)/time.Millisecond), tx.db.addr, tx.db.addr, "tx:Exec")
	if err != nil {
		err = errors.Wrapf(err, "exec:%s, args:%+v", query, args)
	}
	return
}

// Query executes a query that returns rows, typically a SELECT.
func (tx *Tx) Query(query string, args ...interface{}) (rows *Rows, err error) {
	if tx.trace != nil {
		_, span := tx.trace.Start(tx.c, "tx.Query")
		defer span.End()

		span.SetAttributes(
			semconv.DBSystemMySQL,
			attribute.String("db.instance", tx.db.addr),
			semconv.PeerServiceKey.String("mysql"),
			semconv.DBStatementKey.String(fmt.Sprintf("exec: %s, args: %+v", query, args)),
		)
	}

	now := time.Now()
	defer slowLog(fmt.Sprintf("Query query: %s, args: %+v", query, args), now)
	defer func() {
		_metricReqDur.Observe(int64(time.Since(now)/time.Millisecond), tx.db.addr, tx.db.addr, "tx:query")
	}()
	// nolint: rowserrcheck
	rs, err := tx.tx.QueryContext(tx.c, query, args...)
	if err == nil {
		rows = &Rows{Rows: rs}
	} else {
		err = errors.Wrapf(err, "query:%s, args:%+v", query, args)
	}
	return
}

// QueryRow executes a query that is expected to return at most one row.
// QueryRow always returns a non-nil value. Errors are deferred until Row's
// Scan method is called.
func (tx *Tx) QueryRow(query string, args ...interface{}) *Row {
	if tx.trace != nil {
		_, span := tx.trace.Start(tx.c, "tx.QueryRow")
		defer span.End()

		span.SetAttributes(
			semconv.DBSystemMySQL,
			attribute.String("db.instance", tx.db.addr),
			semconv.PeerServiceKey.String("mysql"),
			semconv.DBStatementKey.String(fmt.Sprintf("exec: %s, args: %+v", query, args)),
		)
	}

	now := time.Now()
	defer slowLog(fmt.Sprintf("QueryRow query: %s, args: %+v", query, args), now)
	defer func() {
		_metricReqDur.Observe(int64(time.Since(now)/time.Millisecond), tx.db.addr, tx.db.addr, "tx:QueryRow")
	}()
	r := tx.tx.QueryRowContext(tx.c, query, args...)
	return &Row{Row: r, db: tx.db, query: query, args: args}
}

// Stmt returns a transaction-specific prepared statement from an existing statement.
func (tx *Tx) Stmt(stmt *Stmt) *Stmt {
	as, ok := stmt.stmt.Load().(*sql.Stmt)
	if !ok {
		return nil
	}
	ts := tx.tx.StmtContext(tx.c, as)
	st := &Stmt{query: stmt.query, tx: true, trace: tx.trace, db: tx.db}
	st.stmt.Store(ts)
	return st
}

// Prepare creates a prepared statement for use within a transaction.
// The returned statement operates within the transaction and can no longer be
// used once the transaction has been committed or rolled back.
// To use an existing prepared statement on this transaction, see Tx.Stmt.
func (tx *Tx) Prepare(query string) (*Stmt, error) {
	if tx.trace != nil {
		_, span := tx.trace.Start(tx.c, "tx.Prepare")
		defer span.End()

		span.SetAttributes(
			semconv.DBSystemMySQL,
			attribute.String("db.instance", tx.db.addr),
			semconv.PeerServiceKey.String("mysql"),
			semconv.DBStatementKey.String(fmt.Sprintf("prepare query: %s", query)),
		)
	}

	defer slowLog(fmt.Sprintf("Prepare query: %s", query), time.Now())
	stmt, err := tx.tx.Prepare(query)
	if err != nil {
		err = errors.Wrapf(err, "prepare %s", query)
		return nil, err
	}
	st := &Stmt{query: query, tx: true, trace: tx.trace, db: tx.db}
	st.stmt.Store(stmt)
	return st, nil
}
