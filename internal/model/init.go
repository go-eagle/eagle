package model

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/go-eagle/eagle/pkg/storage/orm"
)

const (
	// DefaultDatabase default database
	DefaultDatabase = "default"
	// UserDatabase user database
	UserDatabase = "user"
)

// Init 初始化数据库
func Init() {
	err := orm.New(
		DefaultDatabase,
		UserDatabase,
	)
	if err != nil {
		panic(fmt.Sprintf("new orm database err: %v", err))
	}
}

// GetDB 返回默认的数据库
func GetDB() (*gorm.DB, error) {
	return orm.GetDB(DefaultDatabase)
}

// GetUserDB 获取用户数据库实例
func GetUserDB() (*gorm.DB, error) {
	return orm.GetDB(UserDatabase)
}
