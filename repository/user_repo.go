package repository

import (
	"github.com/1024casts/snake/model"

	"github.com/jinzhu/gorm"
	"github.com/lexkong/log"
)

// IUserRepo 定义用户仓库接口
type IUserRepo interface {
	CreateUser(db *gorm.DB, user model.UserModel) (id uint64, err error)
	GetUserByID(id uint64) (*model.UserModel, error)
	GetUserByPhone(phone int) (*model.UserModel, error)
	GetUserByEmail(email string) (*model.UserModel, error)
	GetUsersByIds(ids []uint64) ([]*model.UserModel, error)
	Update(userMap map[string]interface{}, id uint64) error
}

// userRepo 用户仓库
type userRepo struct{}

// NewUserRepo 实例化用户仓库
func NewUserRepo() IUserRepo {
	return &userRepo{}
}

// CreateUser 创建用户
func (repo *userRepo) CreateUser(db *gorm.DB, user model.UserModel) (id uint64, err error) {
	err = db.Create(&user).Error
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

// GetUserByID 获取用户
func (repo *userRepo) GetUserByID(id uint64) (*model.UserModel, error) {
	user := &model.UserModel{}
	result := model.GetDB().Where("id = ?", id).First(user)

	return user, result.Error
}

// GetUserByPhone 根据手机号获取用户
func (repo *userRepo) GetUserByPhone(phone int) (*model.UserModel, error) {
	user := model.UserModel{}
	result := model.GetDB().Where("phone = ?", phone).First(&user)

	log.Warnf("select result: %v", user)

	return &user, result.Error
}

// GetUserByEmail 根据邮箱获取手机号
func (repo *userRepo) GetUserByEmail(phone string) (*model.UserModel, error) {
	user := model.UserModel{}
	result := model.GetDB().Where("email = ?", phone).First(&user)

	log.Warnf("select result: %v", user)

	return &user, result.Error
}

// GetUsersByIds 批量获取用户
func (repo *userRepo) GetUsersByIds(ids []uint64) ([]*model.UserModel, error) {
	users := make([]*model.UserModel, 0)
	result := model.GetDB().Where("id in (?)", ids).Find(&users)

	return users, result.Error
}

// Update 更新用户信息
func (repo *userRepo) Update(userMap map[string]interface{}, id uint64) error {
	user, err := repo.GetUserByID(id)
	if err != nil {
		return err
	}

	return model.GetDB().Model(user).Updates(userMap).Error
}
