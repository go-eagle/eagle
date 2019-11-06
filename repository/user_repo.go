package repository

import (
	"github.com/1024casts/snake/model"
	"github.com/pkg/errors"

	"github.com/1024casts/snake/pkg/db"
)

type UserRepository interface {
	Repository
}

type UserRepositoryImpl struct {
	Repository
}

func (repository UserRepositoryImpl) CreateTable(transaction db.Connection) error {
	m := &model.UserModel{}
	tx := transaction.Conn()
	if tx.HasTable(m) {
		return nil
	}
	if err := tx.AutoMigrate(m).Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{&BaseRepository{}}
}
