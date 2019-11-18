package db

import (
	"fmt"
	"strings"

	"github.com/1024casts/snake/pkg/errno"
	"github.com/1024casts/snake/pkg/valid"
	"github.com/pkg/errors"
)

type queryOperator string

const (
	AndQueryOperator queryOperator = "AND"
	OrQueryOperator  queryOperator = "OR"
)

type QueryBuilder interface {
	Or(q Query) Query
	And(q Query) Query
}

type Query interface {
	QueryBuilder
	Where() string
	Values() []interface{}
}

type queryImpl struct {
	where  string
	values []interface{}
}

func (query *queryImpl) combineQueryWithOperator(q Query, operator queryOperator) Query {
	query.where = fmt.Sprintf("(%s) %s (%s)", query.Where(), operator, q.Where())
	query.values = append(query.values, q.Values()...)
	return query
}

func (query *queryImpl) Or(q Query) Query {
	return query.combineQueryWithOperator(q, OrQueryOperator)
}

func (query *queryImpl) And(q Query) Query {
	return query.combineQueryWithOperator(q, AndQueryOperator)
}

func (query *queryImpl) Where() string {
	return query.where
}

func (query *queryImpl) Values() []interface{} {
	return query.values
}

func NewQuery(w string, values ...interface{}) (Query, error) {
	if valid.IsZero(w) {
		return nil, errors.Wrap(errno.ErrParam, "invalid wheres")
	}
	if strings.Count(w, "?") != len(values) {
		return nil, errors.Wrap(errno.ErrParam, "invalid values")
	}
	return &queryImpl{
		where:  w,
		values: values,
	}, nil
}
