package dao

import (
	"context"

	"gorm.io/gorm"

	"github.com/1024casts/snake/internal/cache"
)

var (
	ErrNotFound = gorm.ErrRecordNotFound
)

// Dao mysql struct
type Dao struct {
	db        *gorm.DB
	userCache *cache.Cache
}

// New new a Dao and return
func New(db *gorm.DB) (d *Dao) {
	d = &Dao{
		db:        db,
		userCache: cache.NewUserCache(),
	}
	return
}

// Ping ping mysql
func (d *Dao) Ping(c context.Context) error {
	return nil
}

// Close release mysql connection
func (d *Dao) Close() {

}
