package repository

import (
	"github.com/1024casts/snake/model"
	"github.com/jinzhu/gorm"
	"github.com/lexkong/log"
)

type UserRepo interface {
	CreateUser(db *gorm.DB, user model.UserModel) (id uint64, err error)
	GetUserById(id uint64) (*model.UserModel, error)
	GetUserByPhone(phone int) (*model.UserModel, error)
	GetUsersByIds(ids []uint64) ([]*model.UserModel, error)
	Update(userMap map[string]interface{}, id uint64) error
}

type UserRepoImpl struct {
}

func NewUserRepo() UserRepo {
	return &UserRepoImpl{}
}

func (repo *UserRepoImpl) CreateUser(db *gorm.DB, user model.UserModel) (id uint64, err error) {
	err = db.Create(&user).Error
	if err != nil {
		return 0, err
	}

	return user.Id, nil
}

func (repo *UserRepoImpl) GetUserById(id uint64) (*model.UserModel, error) {
	user := &model.UserModel{}
	result := model.GetDB().Where("id = ?", id).First(user)

	return user, result.Error
}

func (repo *UserRepoImpl) GetUserByPhone(phone int) (*model.UserModel, error) {
	user := model.UserModel{}
	result := model.GetDB().Where("phone = ?", phone).First(&user)

	log.Warnf("select result: %v", user)

	return &user, result.Error
}

func (repo *UserRepoImpl) GetUsersByIds(ids []uint64) ([]*model.UserModel, error) {
	users := make([]*model.UserModel, 0)
	result := model.GetDB().Where("id in (?)", ids).Find(&users)

	return users, result.Error
}

func (repo *UserRepoImpl) Update(userMap map[string]interface{}, id uint64) error {
	user, err := repo.GetUserById(id)
	if err != nil {
		return err
	}

	return model.GetDB().Model(user).Updates(userMap).Error
}
