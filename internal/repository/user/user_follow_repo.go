package user

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"github.com/1024casts/snake/internal/model"
	"github.com/1024casts/snake/pkg/log"
)

// FollowRepo 定义用户仓库接口
type FollowRepo interface {
	CreateUserFollow(db *gorm.DB, userID, followedUID uint64) error
	CreateUserFans(db *gorm.DB, userID, followerUID uint64) error
	UpdateUserFollowStatus(db *gorm.DB, userID, followedUID uint64, status int) error
	UpdateUserFansStatus(db *gorm.DB, userID, followerUID uint64, status int) error
	GetFollowingUserList(userID, lastID uint64, limit int) ([]*model.UserFollowModel, error)
	GetFollowerUserList(userID, lastID uint64, limit int) ([]*model.UserFansModel, error)
	GetFollowByUIds(userID uint64, followingUID []uint64) (map[uint64]*model.UserFollowModel, error)
	GetFansByUIds(userID uint64, followerUID []uint64) (map[uint64]*model.UserFansModel, error)
}

// userFollowRepo 用户仓库
type userFollowRepo struct{}

// NewUserFollowRepo 实例化用户仓库
func NewUserFollowRepo() FollowRepo {
	return &userFollowRepo{}
}

func (repo *userFollowRepo) CreateUserFollow(db *gorm.DB, userID, followedUID uint64) error {
	return db.Exec("insert into user_follow set user_id=?, followed_uid=?, status=1, created_at=? on duplicate key update status=1, updated_at=?",
		userID, followedUID, time.Now(), time.Now()).Error
}

func (repo *userFollowRepo) CreateUserFans(db *gorm.DB, userID, followerUID uint64) error {
	return db.Exec("insert into user_fans set user_id=?, follower_uid=?, status=1, created_at=? on duplicate key update status=1, updated_at=?",
		userID, followerUID, time.Now(), time.Now()).Error
}

func (repo *userFollowRepo) UpdateUserFollowStatus(db *gorm.DB, userID, followedUID uint64, status int) error {
	userFollow := model.UserFollowModel{}
	return db.Model(&userFollow).Where("user_id=? and followed_uid=?", userID, followedUID).
		Updates(map[string]interface{}{"status": status, "updated_at": time.Now()}).Error
}

func (repo *userFollowRepo) UpdateUserFansStatus(db *gorm.DB, userID, followerUID uint64, status int) error {
	userFans := model.UserFansModel{}
	return db.Model(&userFans).Where("user_id=? and follower_uid=?", userID, followerUID).
		Updates(map[string]interface{}{"status": status, "updated_at": time.Now()}).Error
}

func (repo *userFollowRepo) GetFollowingUserList(userID, lastID uint64, limit int) ([]*model.UserFollowModel, error) {
	userFollowList := make([]*model.UserFollowModel, 0)
	db := model.GetDB()
	result := db.Where("user_id=? AND id<=? and status=1", userID, lastID).
		Order("id desc").
		Limit(limit).Find(&userFollowList)

	if err := result.Error; err != nil {
		log.Warnf("[userFollow_service] get user follow list err, %v", err)
		return nil, err
	}

	return userFollowList, nil
}

func (repo *userFollowRepo) GetFollowerUserList(userID, lastID uint64, limit int) ([]*model.UserFansModel, error) {
	userFollowerList := make([]*model.UserFansModel, 0)
	db := model.GetDB()
	result := db.Where("user_id=? AND id<=? and status=1", userID, lastID).
		Order("id desc").
		Limit(limit).Find(&userFollowerList)

	if err := result.Error; err != nil {
		log.Warnf("[userFollow_service] get user follow list err, %v", err)
		return nil, err
	}

	return userFollowerList, nil
}

// 获取自己对关注列表的关注信息
func (repo *userFollowRepo) GetFollowByUIds(userID uint64, followingUID []uint64) (map[uint64]*model.UserFollowModel, error) {
	userFollowModel := make([]*model.UserFollowModel, 0)
	retMap := make(map[uint64]*model.UserFollowModel)

	err := model.GetDB().
		Where("user_id=? AND followed_uid in (?) ", userID, followingUID).
		Find(&userFollowModel).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return retMap, errors.Wrap(err, "[user_follow] get user follow err")
	}

	for _, v := range userFollowModel {
		retMap[v.FollowedUID] = v
	}

	return retMap, nil
}

// 获取自己对关注列表的被关注信息
func (repo *userFollowRepo) GetFansByUIds(userID uint64, followerUID []uint64) (map[uint64]*model.UserFansModel, error) {
	userFansModel := make([]*model.UserFansModel, 0)
	retMap := make(map[uint64]*model.UserFansModel)

	err := model.GetDB().
		Where("user_id=? AND follower_uid in (?) ", userID, followerUID).
		Find(&userFansModel).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return retMap, errors.Wrap(err, "[user_follow] get user fans err")
	}

	for _, v := range userFansModel {
		retMap[v.UserID] = v
	}

	return retMap, nil
}
