package user

import (
	"time"

	"github.com/1024casts/snake/model"
	"github.com/1024casts/snake/pkg/log"
	"github.com/jinzhu/gorm"
)

// Repo 定义用户仓库接口
type FollowRepo interface {
	CreateUserFollow(db *gorm.DB, userId, followedUID uint64) error
	CreateUserFans(db *gorm.DB, userId, followerUID uint64) error
	UpdateUserFollowStatus(db *gorm.DB, userId, followedUId uint64, status int) error
	UpdateUserFansStatus(db *gorm.DB, userId, followerUID uint64, status int) error
	GetFollowingUserList(userId, lastId uint64, limit int) ([]*model.UserFollowModel, error)
	GetFollowerUserList(userId, lastId uint64, limit int) ([]*model.UserFansModel, error)
}

// userFollowRepo 用户仓库
type userFollowRepo struct{}

// NewUserFollowRepo 实例化用户仓库
func NewUserFollowRepo() FollowRepo {
	return &userFollowRepo{}
}

func (repo *userFollowRepo) CreateUserFollow(db *gorm.DB, userId, followedUId uint64) error {
	return db.Exec("insert into user_follow set user_id=?, followed_uid=?, status=1, created_at=? on duplicate key update status=1, updated_at=?",
		userId, followedUId, time.Now(), time.Now()).Error
}

func (repo *userFollowRepo) CreateUserFans(db *gorm.DB, userId, followerUID uint64) error {
	return db.Exec("insert into user_fans set user_id=?, follower_uid=?, status=1, created_at=? on duplicate key update status=1, updated_at=?",
		userId, followerUID, time.Now(), time.Now()).Error
}

func (repo *userFollowRepo) UpdateUserFollowStatus(db *gorm.DB, userId, followedUId uint64, status int) error {
	userFollow := model.UserFollowModel{}
	return db.Model(&userFollow).Where("user_id=? and followed_uid=? and status=?", userId, followedUId, status).
		Updates(map[string]interface{}{"status": 0, "updated_at": time.Now()}).Error
}

func (repo *userFollowRepo) UpdateUserFansStatus(db *gorm.DB, userId, followerUID uint64, status int) error {
	userFans := model.UserFansModel{}
	return db.Model(&userFans).Where("user_id=? and follower_uid=? and status=?", userId, followerUID, status).
		Updates(map[string]interface{}{"status": 0, "updated_at": time.Now()}).Error
}

func (repo *userFollowRepo) GetFollowingUserList(userId, lastId uint64, limit int) ([]*model.UserFollowModel, error) {
	userFollowList := make([]*model.UserFollowModel, 0)
	db := model.GetDB()
	result := db.Where("user_id=? AND id<=? and status=1", userId, lastId).
		Order("id desc").
		Limit(limit).Find(&userFollowList)

	if err := result.Error; err != nil {
		log.Warnf("[userFollow_service] get user follow list err, %v", err)
		return nil, err
	}

	return userFollowList, nil
}

func (repo *userFollowRepo) GetFollowerUserList(userId, lastId uint64, limit int) ([]*model.UserFansModel, error) {
	userFollowerList := make([]*model.UserFansModel, 0)
	db := model.GetDB()
	result := db.Where("user_id=? AND id<=? and status=1", userId, lastId).
		Order("id desc").
		Limit(limit).Find(&userFollowerList)

	if err := result.Error; err != nil {
		log.Warnf("[userFollow_service] get user follow list err, %v", err)
		return nil, err
	}

	return userFollowerList, nil
}
