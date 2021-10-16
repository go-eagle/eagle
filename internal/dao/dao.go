package dao

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"

	"github.com/go-eagle/eagle/internal/cache"
	"github.com/go-eagle/eagle/pkg/conf"
	"github.com/go-eagle/eagle/pkg/storage/sql"
)

var (
	// ErrNotFound data is not exist
	ErrNotFound = gorm.ErrRecordNotFound
)

// Dao mysql struct
type Dao struct {
	orm       *gorm.DB
	db        *sql.DB
	tracer    trace.Tracer
	userCache *cache.Cache
}

// New new a Dao and return
func New(cfg *conf.Config, db *gorm.DB) *Dao {
	return &Dao{
		orm:       db,
		db:        sql.NewMySQL(cfg.MySQL),
		tracer:    otel.Tracer("dao"),
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
