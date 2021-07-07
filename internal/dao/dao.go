package dao

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"gorm.io/gorm"

	"github.com/1024casts/snake/internal/cache"
)

var (
	ErrNotFound = gorm.ErrRecordNotFound
)

// Dao mysql struct
type Dao struct {
	db        *gorm.DB
	tracer    trace.Tracer
	userCache *cache.Cache
}

// New new a Dao and return
func New(db *gorm.DB) *Dao {
	return &Dao{
		db:        db,
		tracer:    otel.GetTracerProvider().Tracer("dao"),
		userCache: cache.NewUserCache(),
	}
}

// Ping ping mysql
func (d *Dao) Ping(c context.Context) error {
	return nil
}

// Close release mysql connection
func (d *Dao) Close() {

}
