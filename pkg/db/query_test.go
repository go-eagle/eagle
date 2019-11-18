package db

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryImpl_combineQueryWithOperator(t *testing.T) {
	type fields struct {
		where  string
		values []interface{}
	}
	type args struct {
		q        Query
		operator queryOperator
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Query
	}{
		{
			name: "and",
			fields: fields{
				where:  "name=?",
				values: []interface{}{"sangil"},
			},
			args: args{
				q: &queryImpl{
					where:  "id=?",
					values: []interface{}{1},
				},
				operator: AndQueryOperator,
			},
			want: &queryImpl{
				where:  "(name=?) AND (id=?)",
				values: []interface{}{"sangil", 1},
			},
		},
		{
			name: "or",
			fields: fields{
				where:  "name=?",
				values: []interface{}{"sangil"},
			},
			args: args{
				q: &queryImpl{
					where:  "id=?",
					values: []interface{}{1},
				},
				operator: OrQueryOperator,
			},
			want: &queryImpl{
				where:  "(name=?) OR (id=?)",
				values: []interface{}{"sangil", 1},
			},
		},
		{
			name: "or",
			fields: fields{
				where:  "name=? AND id=?",
				values: []interface{}{"sangil", 1},
			},
			args: args{
				q: &queryImpl{
					where:  "age=?",
					values: []interface{}{1},
				},
				operator: OrQueryOperator,
			},
			want: &queryImpl{
				where:  "(name=? AND id=?) OR (age=?)",
				values: []interface{}{"sangil", 1, 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := &queryImpl{
				where:  tt.fields.where,
				values: tt.fields.values,
			}
			got := query.combineQueryWithOperator(tt.args.q, tt.args.operator)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestQueryImpl_Or(t *testing.T) {
	type fields struct {
		where  string
		values []interface{}
	}
	type args struct {
		q Query
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Query
	}{
		{
			name: "and",
			fields: fields{
				where:  "name=?",
				values: []interface{}{"sangil"},
			},
			args: args{
				q: &queryImpl{
					where:  "id=?",
					values: []interface{}{1},
				},
			},
			want: &queryImpl{
				where:  "(name=?) OR (id=?)",
				values: []interface{}{"sangil", 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := &queryImpl{
				where:  tt.fields.where,
				values: tt.fields.values,
			}
			if got := query.Or(tt.args.q); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QueryImpl.Or() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueryImpl_And(t *testing.T) {
	type fields struct {
		where  string
		values []interface{}
	}
	type args struct {
		q Query
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Query
	}{
		{
			name: "and",
			fields: fields{
				where:  "name=?",
				values: []interface{}{"sangil"},
			},
			args: args{
				q: &queryImpl{
					where:  "id=?",
					values: []interface{}{1},
				},
			},
			want: &queryImpl{
				where:  "(name=?) AND (id=?)",
				values: []interface{}{"sangil", 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := &queryImpl{
				where:  tt.fields.where,
				values: tt.fields.values,
			}
			if got := query.And(tt.args.q); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QueryImpl.And() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewQuery(t *testing.T) {
	type args struct {
		w      string
		values []interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    Query
		wantErr bool
	}{
		{
			name: "pass",
			args: args{
				w:      "id=?",
				values: []interface{}{1},
			},
			want: &queryImpl{
				where:  "id=?",
				values: []interface{}{1},
			},
			wantErr: false,
		},
		{
			name: "invalid where",
			args: args{
				w:      "",
				values: []interface{}{1},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid values",
			args: args{
				w:      "id=1",
				values: []interface{}{1},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewQuery(tt.args.w, tt.args.values...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
