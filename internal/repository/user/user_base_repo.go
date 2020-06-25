package user

import (
	"fmt"
	"time"

	"github.com/1024casts/snake/pkg/redis"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"github.com/1024casts/snake/internal/cache/user"
	"github.com/1024casts/snake/internal/model"
	"github.com/1024casts/snake/pkg/log"
)

// Repo 定义用户仓库接口
type BaseRepo interface {
	Create(db *gorm.DB, user model.UserBaseModel) (id uint64, err error)
	Update(db *gorm.DB, id uint64, userMap map[string]interface{}) error
	GetUserByID(db *gorm.DB, id uint64) (*model.UserBaseModel, error)
	GetUsersByIds(db *gorm.DB, ids []uint64) ([]*model.UserBaseModel, error)
	GetUserByPhone(db *gorm.DB, phone int) (*model.UserBaseModel, error)
	GetUserByEmail(db *gorm.DB, email string) (*model.UserBaseModel, error)
}

// userRepo 用户仓库
type userRepo struct {
	userCache *user.Cache
}

// NewUserRepo 实例化用户仓库
func NewUserRepo() BaseRepo {
	return &userRepo{
		userCache: user.NewUserCache(),
	}
}

// Create 创建用户
func (repo *userRepo) Create(db *gorm.DB, user model.UserBaseModel) (id uint64, err error) {
	err = db.Create(&user).Error
	if err != nil {
		return 0, errors.Wrap(err, "[user_repo] create user err")
	}

	return user.ID, nil
}

// Update 更新用户信息
func (repo *userRepo) Update(db *gorm.DB, id uint64, userMap map[string]interface{}) error {
	user, err := repo.GetUserByID(db, id)
	if err != nil {
		return errors.Wrap(err, "[user_repo] update user data err")
	}

	// 删除cache
	err = repo.userCache.DelUserBaseCache(id)
	if err != nil {
		log.Warnf("[user_repo] delete user cache err: %v", err)
	}

	return db.Model(user).Updates(userMap).Error
}

// GetUserByID 获取用户
func (repo *userRepo) GetUserByID(db *gorm.DB, id uint64) (*model.UserBaseModel, error) {
	// 从cache获取
	userModel, err := repo.userCache.GetUserBaseCache(id)
	if err != nil {
		return nil, errors.Wrap(err, "[user_repo] get user cache data err")
	}
	if userModel != nil && userModel.ID > 0 {
		return userModel, nil
	}

	// 加锁，防止缓存击穿
	key := fmt.Sprintf("uid:%d", id)
	lock := redis.NewLock(redis.Client, key, 3*time.Second)
	token := lock.GenToken()

	isLock, err := lock.Lock(token)
	if err != nil || !isLock {
		return nil, errors.Wrap(err, "[user_repo] lock err")
	}
	defer lock.Unlock(token)

	data := &model.UserBaseModel{}
	if isLock {
		// 从数据库中获取
		err = db.Where(&model.UserBaseModel{ID: id}).First(data).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, errors.Wrap(err, "[user_repo] get user data err")
		}

		// 写入cache
		err = repo.userCache.SetUserBaseCache(id, data)
		if err != nil {
			return data, errors.Wrap(err, "[user_repo] set user data err")
		}
	}
	return data, nil
}

// GetUsersByIds 批量获取用户
func (repo *userRepo) GetUsersByIds(db *gorm.DB, userIDs []uint64) ([]*model.UserBaseModel, error) {
	users := make([]*model.UserBaseModel, 0)

	// 从cache批量获取
	userCacheMap, err := repo.userCache.MultiGetUserBaseCache(userIDs)
	if err != nil {
		return users, errors.Wrap(err, "[user_repo] multi get user cache data err")
	}

	// 查询未命中
	for _, userID := range userIDs {
		idx := repo.userCache.GetUserBaseCacheKey(userID)
		userModel, ok := userCacheMap[idx]
		if !ok {
			userModel, err = repo.GetUserByID(db, userID)
			if err != nil {
				log.Warnf("get user model err: %v", err)
				continue
			}
		}
		users = append(users, userModel)
	}
	return users, nil
}

// GetUserByPhone 根据手机号获取用户
func (repo *userRepo) GetUserByPhone(db *gorm.DB, phone int) (*model.UserBaseModel, error) {
	user := model.UserBaseModel{}
	err := db.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, errors.Wrap(err, "[user_repo] get user err by phone")
	}

	return &user, nil
}

// GetUserByEmail 根据邮箱获取手机号
func (repo *userRepo) GetUserByEmail(db *gorm.DB, phone string) (*model.UserBaseModel, error) {
	user := model.UserBaseModel{}
	err := db.Where("email = ?", phone).First(&user).Error
	if err != nil {
		return nil, errors.Wrap(err, "[user_repo] get user err by email")
	}

	return &user, nil
}
