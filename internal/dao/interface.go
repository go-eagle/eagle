package dao

import (
	"context"

	"gorm.io/gorm"

	"github.com/go-eagle/eagle/internal/model"
)

var _ dao = (*Dao)(nil)

// dao 定义用户仓库接口
type dao interface {
	// BaseUser
	CreateUser(ctx context.Context, user *model.UserBaseModel) (id uint64, err error)
	UpdateUser(ctx context.Context, id uint64, userMap map[string]interface{}) error
	GetOneUser(ctx context.Context, id uint64) (*model.UserBaseModel, error)
	GetUsersByIds(ctx context.Context, ids []uint64) ([]*model.UserBaseModel, error)
	GetUserByPhone(ctx context.Context, phone int64) (*model.UserBaseModel, error)
	GetUserByEmail(ctx context.Context, email string) (*model.UserBaseModel, error)

	// Follow
	CreateUserFollow(ctx context.Context, db *gorm.DB, userID, followedUID uint64) error
	CreateUserFans(ctx context.Context, db *gorm.DB, userID, followerUID uint64) error
	UpdateUserFollowStatus(ctx context.Context, db *gorm.DB, userID, followedUID uint64, status int) error
	UpdateUserFansStatus(ctx context.Context, db *gorm.DB, userID, followerUID uint64, status int) error
	GetFollowingUserList(ctx context.Context, userID, lastID uint64, limit int) ([]*model.UserFollowModel, error)
	GetFollowerUserList(ctx context.Context, userID, lastID uint64, limit int) ([]*model.UserFansModel, error)
	GetFollowByUIds(ctx context.Context, userID uint64, followingUID []uint64) (map[uint64]*model.UserFollowModel, error)
	GetFansByUIds(ctx context.Context, userID uint64, followerUID []uint64) (map[uint64]*model.UserFansModel, error)

	// stat
	IncrFollowCount(ctx context.Context, db *gorm.DB, userID uint64, step int) error
	IncrFollowerCount(ctx context.Context, db *gorm.DB, userID uint64, step int) error
	GetUserStatByID(ctx context.Context, userID uint64) (*model.UserStatModel, error)
	GetUserStatByIDs(ctx context.Context, userID []uint64) (map[uint64]*model.UserStatModel, error)

	Close()
}
