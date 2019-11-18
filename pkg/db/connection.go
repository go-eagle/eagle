package db

import (
	"fmt"
	"time"

	"github.com/spf13/viper"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/realsangil/apimonitor/pkg/rserrors"
)

const ErrInvalidTransaction = rserrors.Error("invalid transaction")

var (
	conn *defaultConnection
)

func GetConnection() Connection {
	return conn
}

type Connection interface {
	Conn() *gorm.DB
	Begin() (Connection, error)
	Close() error
	Rollback() error
	Commit() error
}

type defaultConnection struct {
	db *gorm.DB
}

func NewConnection(tx *gorm.DB) Connection {
	return &defaultConnection{tx}
}

func (conn *defaultConnection) Conn() *gorm.DB {
	return conn.db
}

func (conn *defaultConnection) Begin() (Connection, error) {
	if conn.db == nil {
		return nil, errors.New("connection is nil")
	}
	tx := conn.db.Begin().Set("gorm:table_options", "ENGINE=InnoDB charset=utf8mb4")
	return &defaultConnection{db: tx}, nil
}

func (conn *defaultConnection) Close() error {
	if conn.db == nil {
		return errors.New("connection is nil")
	}
	return conn.db.Close()
}

func (conn *defaultConnection) Rollback() error {
	return conn.Conn().Rollback().Error
}

func (conn *defaultConnection) Commit() error {
	return conn.Conn().Commit().Error
}

func Init(config DatabaseConfig) error {
	cfg := fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=true&loc=%s",
		config.GetUsername(),
		config.GetPassword(),
		config.GetHost(),
		config.GetPort(),
		config.GetDatabaseName(),
		"Local",
	)
	db, err := gorm.Open("mysql", cfg)

	// 如果没有连接上，则尝试创建数据库
	if err != nil {
		db, err = gorm.Open("mysql", fmt.Sprintf(
			"%v:%v@tcp(%v:%v)/?charset=utf8mb4&parseTime=true&sql_mode=STRICT_ALL_TABLES",
			config.GetUsername(),
			config.GetPassword(),
			config.GetHost(),
			config.GetPort(),
		))
		if err != nil {
			return errors.WithStack(err)
		}

		db = db.Set("gorm:table_options", "ENGINE=InnoDB charset=utf8mb4")
		_, err = db.DB().Exec("CREATE DATABASE " + config.GetDatabaseName())
		if err != nil {
			return errors.WithStack(err)
		}

		if err = db.Close(); err != nil {
			return errors.WithStack(err)
		}

		db, err = gorm.Open("mysql", cfg)
		if err != nil {
			return errors.WithStack(err)
		}

		if err = db.DB().Ping(); err != nil {
			return errors.WithStack(err)
		}
	}
	db = db.Set("gorm:table_options", "ENGINE=InnoDB charset=utf8mb4")

	db.DB().SetMaxOpenConns(viper.GetInt("grom.max_open_conn")) // 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	db.DB().SetMaxIdleConns(viper.GetInt("grom.max_idle_conn")) // 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
	db.DB().SetConnMaxLifetime(time.Minute * viper.GetDuration("grom.conn_max_life_time"))

	if config.GetVerbose() {
		db = db.Debug()
	}

	conn = &defaultConnection{db: db}
	return nil
}
