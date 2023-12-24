package repository

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/go-eagle/eagle/internal/model"
	"github.com/go-eagle/eagle/pkg/log"
)

// CreateUserFollow .
func (d *repository) CreateUserFollow(ctx context.Context, db *gorm.DB, userID, followedUID uint64) error {
	return db.Exec("insert into user_follow set user_id=?, followed_uid=?, status=1, created_at=? on duplicate key update status=1, updated_at=?",
		userID, followedUID, time.Now(), time.Now()).Error
}

// CreateUserFans .
func (d *repository) CreateUserFans(ctx context.Context, db *gorm.DB, userID, followerUID uint64) error {
	return db.Exec("insert into user_fans set user_id=?, follower_uid=?, status=1, created_at=? on duplicate key update status=1, updated_at=?",
		userID, followerUID, time.Now(), time.Now()).Error
}

// UpdateUserFollowStatus .
func (d *repository) UpdateUserFollowStatus(ctx context.Context, db *gorm.DB, userID, followedUID uint64, status int) error {
	userFollow := model.UserFollowModel{}
	return db.Model(&userFollow).Where("user_id=? and followed_uid=?", userID, followedUID).
		Updates(map[string]interface{}{"status": status, "updated_at": time.Now()}).Error
}

// UpdateUserFansStatus .
func (d *repository) UpdateUserFansStatus(ctx context.Context, db *gorm.DB, userID, followerUID uint64, status int) error {
	userFans := model.UserFansModel{}
	return db.Model(&userFans).Where("user_id=? and follower_uid=?", userID, followerUID).
		Updates(map[string]interface{}{"status": status, "updated_at": time.Now()}).Error
}

// GetFollowingUserList .
func (d *repository) GetFollowingUserList(ctx context.Context, userID, lastID uint64, limit int) ([]*model.UserFollowModel, error) {
	userFollowList := make([]*model.UserFollowModel, 0)
	db, _ := model.GetDB()
	result := db.Where("user_id=? AND id<=? and status=1", userID, lastID).
		Order("id desc").
		Limit(limit).Find(&userFollowList)

	if err := result.Error; err != nil {
		log.Warnf("[userFollow_service] get user follow list err, %v", err)
		return nil, err
	}

	return userFollowList, nil
}

// GetFollowerUserList get follower user list
func (d *repository) GetFollowerUserList(ctx context.Context, userID, lastID uint64, limit int) ([]*model.UserFansModel, error) {
	userFollowerList := make([]*model.UserFansModel, 0)
	db, _ := model.GetDB()
	result := db.Where("user_id=? AND id<=? and status=1", userID, lastID).
		Order("id desc").
		Limit(limit).Find(&userFollowerList)

	if err := result.Error; err != nil {
		log.Warnf("[userFollow_service] get user follow list err, %v", err)
		return nil, err
	}

	return userFollowerList, nil
}

// GetFollowByUIds 获取自己对关注列表的关注信息
func (d *repository) GetFollowByUIds(ctx context.Context, userID uint64, followingUID []uint64) (map[uint64]*model.UserFollowModel, error) {
	userFollowModel := make([]*model.UserFollowModel, 0)
	retMap := make(map[uint64]*model.UserFollowModel)
	db, _ := model.GetDB()
	err := db.
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

// GetFansByUIds 获取自己对关注列表的被关注信息
func (d *repository) GetFansByUIds(ctx context.Context, userID uint64, followerUID []uint64) (map[uint64]*model.UserFansModel, error) {
	userFansModel := make([]*model.UserFansModel, 0)
	retMap := make(map[uint64]*model.UserFansModel)
	db, _ := model.GetDB()
	err := db.
		Where("user_id=? AND follower_uid in (?) ", userID, followerUID).
		Find(&userFansModel).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return retMap, errors.Wrap(err, "[user_follow] get user fans err")
	}

	for _, v := range userFansModel {
		retMap[v.UserID] = v
	}

	return retMap, nil
}
