package user

import (
	"time"

	"github.com/jinzhu/gorm"

	"github.com/1024casts/snake/model"
)

// Repo 定义用户仓库接口
type Repo interface {
	Create(db *gorm.DB, user model.UserModel) (id uint64, err error)
	Update(db *gorm.DB, id uint64, userMap map[string]interface{}) error
	GetUserByID(db *gorm.DB, id uint64) (*model.UserModel, error)
	GetUserByPhone(db *gorm.DB, phone int) (*model.UserModel, error)
	GetUserByEmail(db *gorm.DB, email string) (*model.UserModel, error)
	GetUsersByIds(ids []uint64) ([]*model.UserModel, error)
	IncrFollowCount(db *gorm.DB, userID uint64, step int) error
	IncrFollowerCount(db *gorm.DB, userID uint64, step int) error
}

// userRepo 用户仓库
type userRepo struct{}

// NewUserRepo 实例化用户仓库
func NewUserRepo() Repo {
	return &userRepo{}
}

// Create 创建用户
func (repo *userRepo) Create(db *gorm.DB, user model.UserModel) (id uint64, err error) {
	err = db.Create(&user).Error
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

// Update 更新用户信息
func (repo *userRepo) Update(db *gorm.DB, id uint64, userMap map[string]interface{}) error {
	user, err := repo.GetUserByID(db, id)
	if err != nil {
		return err
	}

	return db.Model(user).Updates(userMap).Error
}

// GetUserByID 获取用户
func (repo *userRepo) GetUserByID(db *gorm.DB, id uint64) (*model.UserModel, error) {
	user := &model.UserModel{}
	result := db.Where("id = ?", id).First(user)

	return user, result.Error
}

// GetUserByPhone 根据手机号获取用户
func (repo *userRepo) GetUserByPhone(db *gorm.DB, phone int) (*model.UserModel, error) {
	user := model.UserModel{}
	result := db.Where("phone = ?", phone).First(&user)

	return &user, result.Error
}

// GetUserByEmail 根据邮箱获取手机号
func (repo *userRepo) GetUserByEmail(db *gorm.DB, phone string) (*model.UserModel, error) {
	user := model.UserModel{}
	result := db.Where("email = ?", phone).First(&user)

	return &user, result.Error
}

// GetUsersByIds 批量获取用户
func (repo *userRepo) GetUsersByIds(ids []uint64) ([]*model.UserModel, error) {
	users := make([]*model.UserModel, 0)
	result := model.GetDB().Where("id in (?)", ids).Find(&users)

	return users, result.Error
}

// IncrFollowCount 增加关注数
func (repo *userRepo) IncrFollowCount(db *gorm.DB, userID uint64, step int) error {
	return db.Exec("insert into user_stat set user_id=?, follow_count=1, created_at=? on duplicate key update "+
		"follow_count=follow_count+?, updated_at=?",
		userID, time.Now(), step, time.Now()).Error
}

// IncrFollowerCount 增加粉丝数
func (repo *userRepo) IncrFollowerCount(db *gorm.DB, userID uint64, step int) error {
	return db.Exec("insert into user_stat set user_id=?, follower_count=1, created_at=? on duplicate key update "+
		"follower_count=follower_count+?, updated_at=?",
		userID, time.Now(), step, time.Now()).Error
}
