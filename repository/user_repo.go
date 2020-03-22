package repository

import (
	"github.com/1024casts/snake/model"
	"github.com/jinzhu/gorm"
	"github.com/lexkong/log"
)

type IUserRepo interface {
	CreateUser(db *gorm.DB, user model.UserModel) (id uint64, err error)
	GetUserById(id uint64) (*model.UserModel, error)
	GetUserByPhone(phone int) (*model.UserModel, error)
	GetUserByEmail(email string) (*model.UserModel, error)
	GetUsersByIds(ids []uint64) ([]*model.UserModel, error)
	Update(userMap map[string]interface{}, id uint64) error
}

type UserRepo struct {
}

func NewUserRepo() IUserRepo {
	return &UserRepo{}
}

func (repo *UserRepo) CreateUser(db *gorm.DB, user model.UserModel) (id uint64, err error) {
	err = db.Create(&user).Error
	if err != nil {
		return 0, err
	}

	return user.Id, nil
}

func (repo *UserRepo) GetUserById(id uint64) (*model.UserModel, error) {
	user := &model.UserModel{}
	result := model.GetDB().Where("id = ?", id).First(user)

	return user, result.Error
}

func (repo *UserRepo) GetUserByPhone(phone int) (*model.UserModel, error) {
	user := model.UserModel{}
	result := model.GetDB().Where("phone = ?", phone).First(&user)

	log.Warnf("select result: %v", user)

	return &user, result.Error
}

func (repo *UserRepo) GetUserByEmail(phone string) (*model.UserModel, error) {
	user := model.UserModel{}
	result := model.GetDB().Where("email = ?", phone).First(&user)

	log.Warnf("select result: %v", user)

	return &user, result.Error
}

func (repo *UserRepo) GetUsersByIds(ids []uint64) ([]*model.UserModel, error) {
	users := make([]*model.UserModel, 0)
	result := model.GetDB().Where("id in (?)", ids).Find(&users)

	return users, result.Error
}

func (repo *UserRepo) Update(userMap map[string]interface{}, id uint64) error {
	user, err := repo.GetUserById(id)
	if err != nil {
		return err
	}

	return model.GetDB().Model(user).Updates(userMap).Error
}
