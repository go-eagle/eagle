package user

import (
	"context"
	"fmt"
	"time"

	"github.com/1024casts/snake/pkg/redis"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"github.com/1024casts/snake/internal/cache/user"
	"github.com/1024casts/snake/internal/model"
	"github.com/1024casts/snake/pkg/log"
)

// BaseRepo 定义用户仓库接口
type BaseRepo interface {
	Create(ctx context.Context, user model.UserBaseModel) (id uint64, err error)
	Update(ctx context.Context, id uint64, userMap map[string]interface{}) error
	GetUserByID(ctx context.Context, id uint64) (*model.UserBaseModel, error)
	GetUsersByIds(ctx context.Context, ids []uint64) ([]*model.UserBaseModel, error)
	GetUserByPhone(ctx context.Context, phone int64) (*model.UserBaseModel, error)
	GetUserByEmail(ctx context.Context, email string) (*model.UserBaseModel, error)
	Close()
}

// userBaseRepo 用户仓库
type userBaseRepo struct {
	db        *gorm.DB
	userCache *user.Cache
}

// NewUserRepo 实例化用户仓库
func NewUserRepo(db *gorm.DB) BaseRepo {
	return &userBaseRepo{
		db:        db,
		userCache: user.NewUserCache(),
	}
}

// Create 创建用户
func (repo *userBaseRepo) Create(ctx context.Context, user model.UserBaseModel) (id uint64, err error) {
	err = repo.db.Create(&user).Error
	if err != nil {
		return 0, errors.Wrap(err, "[user_repo] create user err")
	}

	return user.ID, nil
}

// Update 更新用户信息
func (repo *userBaseRepo) Update(ctx context.Context, id uint64, userMap map[string]interface{}) error {
	user, err := repo.GetUserByID(ctx, id)
	if err != nil {
		return errors.Wrap(err, "[user_repo] update user data err")
	}

	// 删除cache
	err = repo.userCache.DelUserBaseCache(id)
	if err != nil {
		log.Warnf("[user_repo] delete user cache err: %v", err)
	}

	return repo.db.Model(user).Updates(userMap).Error
}

// GetUserByID 获取用户
// 缓存的更新策略使用 Cache Aside Pattern
// see: https://coolshell.cn/articles/17416.html
func (repo *userBaseRepo) GetUserByID(ctx context.Context, uid uint64) (userBase *model.UserBaseModel, err error) {
	//var userBase *model.UserBaseModel
	start := time.Now()
	defer func() {
		log.Infof("[repo.user_base] get user by uid: %d cost: %d ns", uid, time.Now().Sub(start).Nanoseconds())
	}()
	// 从cache获取
	userBase, err = repo.userCache.GetUserBaseCache(uid)
	if err != nil {
		// if cache error return, don't request to db
		log.Warnf("[repo.user_base] get user by uid err: %v, uid: %d", err, uid)
		return
	}
	// hit cache
	if userBase != nil {
		log.Infof("[repo.user_base] get user base data from cache, uid: %d", uid)
		return
	}

	// 加锁，防止缓存击穿
	key := fmt.Sprintf("uid:%d", uid)
	lock := redis.NewLock(redis.RedisClient, key, 3*time.Second)
	token := lock.GenToken()

	isLock, err := lock.Lock(token)
	if err != nil || !isLock {
		return nil, errors.Wrapf(err, "[repo.user_base] lock err, key: %s", key)
	}
	defer lock.Unlock(token)

	data := new(model.UserBaseModel)
	if isLock {
		// 从数据库中获取
		err = repo.db.First(data, uid).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, errors.Wrap(err, "[repo.user_base] get user data err")
		}

		// set cache
		err = repo.userCache.SetUserBaseCache(uid, data)
		if err != nil {
			return data, errors.Wrap(err, "[repo.user_base] set user data err")
		}
	}
	return data, nil
}

// GetUsersByIds 批量获取用户
func (repo *userBaseRepo) GetUsersByIds(ctx context.Context, userIDs []uint64) ([]*model.UserBaseModel, error) {
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
			userModel, err = repo.GetUserByID(ctx, userID)
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
func (repo *userBaseRepo) GetUserByPhone(ctx context.Context, phone int64) (*model.UserBaseModel, error) {
	user := model.UserBaseModel{}
	err := repo.db.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, errors.Wrap(err, "[user_repo] get user err by phone")
	}

	return &user, nil
}

// GetUserByEmail 根据邮箱获取手机号
func (repo *userBaseRepo) GetUserByEmail(ctx context.Context, phone string) (*model.UserBaseModel, error) {
	user := model.UserBaseModel{}
	err := repo.db.Where("email = ?", phone).First(&user).Error
	if err != nil {
		return nil, errors.Wrap(err, "[user_repo] get user err by email")
	}

	return &user, nil
}

// Close close db
func (repo *userBaseRepo) Close() {
	repo.db.Close()
}
