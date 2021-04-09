package dao

import (
	"context"
	"fmt"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"

	"github.com/1024casts/snake/internal/model"
	"github.com/1024casts/snake/pkg/cache"
	"github.com/1024casts/snake/pkg/log"
	"github.com/1024casts/snake/pkg/metric/prom"
	"github.com/1024casts/snake/pkg/redis"
)

// Create 创建用户
func (d *Dao) CreateUser(ctx context.Context, user model.UserBaseModel) (id uint64, err error) {
	err = d.db.Create(&user).Error
	if err != nil {
		prom.BusinessErrCount.Incr("mysql: CreateUser")
		return 0, errors.Wrap(err, "[repo.user_base] create user err")
	}

	return user.ID, nil
}

// Update 更新用户信息
func (d *Dao) UpdateUser(ctx context.Context, id uint64, userMap map[string]interface{}) error {
	user, err := d.GetOneUser(ctx, id)
	if err != nil {
		prom.BusinessErrCount.Incr("mysql: getOneUser")
		return errors.Wrap(err, "[repo.user_base] update user data err")
	}

	// 删除cache
	err = d.userCache.DelUserBaseCache(ctx, id)
	if err != nil {
		log.Warnf("[repo.user_base] delete user cache err: %v", err)
	}

	err = d.db.Model(user).Updates(userMap).Error
	if err != nil {
		prom.BusinessErrCount.Incr("mysql: UpdateUser")
	}
	return err
}

// GetUserByID 获取用户
// 缓存的更新策略使用 Cache Aside Pattern
// see: https://coolshell.cn/articles/17416.html
func (d *Dao) GetOneUser(ctx context.Context, uid uint64) (userBase *model.UserBaseModel, err error) {
	// add tracing
	span, ctx := opentracing.StartSpanFromContext(ctx, "userBaseDao.GetOneUser")
	defer span.Finish()

	//var userBase *model.UserBaseModel
	start := time.Now()
	defer func() {
		log.Infof("[repo.user_base] get user by uid: %d cost: %d μs", uid, time.Since(start).Microseconds())
	}()
	// 从cache获取
	userBase, err = d.userCache.GetUserBaseCache(ctx, uid)
	if err != nil {
		if err == cache.ErrPlaceholder {
			return nil, ErrNotFound
		} else if !errors.Is(err, redis.ErrRedisNotFound) {
			// fail fast, if cache error return, don't request to db
			return nil, errors.Wrapf(err, "[repo.user_base] get user by uid: %d", uid)
		}
	}
	// cache hit
	if userBase != nil {
		prom.CacheHit.Incr("getOneUser")
		log.Infof("[repo.user_base] get user base data from cache, uid: %d", uid)
		return
	}

	// use sync/singleflight mode to get data
	// demo see: https://github.com/go-demo/singleflight-demo/blob/master/main.go
	// https://juejin.cn/post/6844904084445593613
	getDataFn := func() (interface{}, error) {
		data := new(model.UserBaseModel)
		// 从数据库中获取
		err = d.db.First(data, uid).Error
		// if data is empty, set not found cache to prevent cache penetration(缓存穿透)
		if errors.Is(err, ErrNotFound) {
			err = d.userCache.SetCacheWithNotFound(ctx, uid)
			if err != nil {
				log.Warnf("[repo.user_base] SetCacheWithNotFound err, uid: %d", uid)
			}
			return nil, ErrNotFound
		} else if err != nil {
			prom.BusinessErrCount.Incr("mysql: getOneUser")
			return nil, errors.Wrapf(err, "[repo.user_base] query db err")
		}

		// set cache
		err = d.userCache.SetUserBaseCache(ctx, uid, data, cache.DefaultExpireTime)
		if err != nil {
			return data, errors.Wrap(err, "[repo.user_base] SetUserBaseCache err")
		}
		return data, nil
	}

	g := singleflight.Group{}
	doKey := fmt.Sprintf("get_user_base_%d", uid)
	val, err, _ := g.Do(doKey, getDataFn)
	if err != nil {
		return nil, errors.Wrap(err, "[repo.user_base] get user base err via single flight do")
	}
	data := val.(*model.UserBaseModel)

	// cache miss
	prom.CacheMiss.Incr("getOneUser")

	return data, nil
}

// GetUsersByIds 批量获取用户
func (d *Dao) GetUsersByIds(ctx context.Context, userIDs []uint64) ([]*model.UserBaseModel, error) {
	users := make([]*model.UserBaseModel, 0)

	// 从cache批量获取
	userCacheMap, err := d.userCache.MultiGetUserBaseCache(ctx, userIDs)
	if err != nil {
		return users, errors.Wrap(err, "[repo.user_base] multi get user cache data err")
	}

	// 查询未命中
	for _, userID := range userIDs {
		idx := d.userCache.GetUserBaseCacheKey(userID)
		userModel, ok := userCacheMap[idx]
		if !ok {
			userModel, err = d.GetOneUser(ctx, userID)
			if err != nil {
				log.Warnf("[repo.user_base] get user model err: %v", err)
				continue
			}
		}
		users = append(users, userModel)
	}
	return users, nil
}

// GetUserByPhone 根据手机号获取用户
func (d *Dao) GetUserByPhone(ctx context.Context, phone int64) (*model.UserBaseModel, error) {
	user := model.UserBaseModel{}
	err := d.db.Where("phone = ?", phone).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.Wrap(err, "[repo.user_base] get user err by phone")
	}

	return &user, nil
}

// GetUserByEmail 根据邮箱获取手机号
func (d *Dao) GetUserByEmail(ctx context.Context, email string) (*model.UserBaseModel, error) {
	userBase := model.UserBaseModel{}
	err := d.db.Where("email = ?", email).First(&userBase).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.Wrap(err, "[repo.user_base] get user err by email")
	}

	return &userBase, nil
}
