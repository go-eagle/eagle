package sql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/trace"
)

// Row row.
type Row struct {
	err error
	*sql.Row
	db     *conn
	query  string
	args   []interface{}
	trace  trace.Tracer
	cancel func()
}

// Scan copies the columns from the matched row into the values pointed at by dest.
func (r *Row) Scan(dest ...interface{}) (err error) {
	defer slowLog(fmt.Sprintf("Scan query: %s args: %+v", r.query, r.args), time.Now())
	if r.trace != nil {
		//r.trace.
		//defer r.trace.Finish(&err)
	}
	if r.err != nil {
		err = r.err
	} else if r.Row == nil {
		err = ErrStmtNil
	}
	if err != nil {
		return
	}
	err = r.Row.Scan(dest...)
	if r.cancel != nil {
		r.cancel()
	}
	r.db.onBreaker(&err)
	if err != ErrNoRows {
		err = errors.Wrapf(err, "query: %s, args: %+v", r.query, r.args)
	}
	return
}

// Rows rows.
type Rows struct {
	*sql.Rows
	cancel func()
}

// Close closes the Rows, preventing further enumeration. If Next is called
// and returns false and there are no further result sets,
// the Rows are closed automatically and it will suffice to check the
// result of Err. Close is idempotent and does not affect the result of Err.
func (rs *Rows) Close() (err error) {
	err = errors.WithStack(rs.Rows.Close())
	if rs.cancel != nil {
		rs.cancel()
	}
	return
}
