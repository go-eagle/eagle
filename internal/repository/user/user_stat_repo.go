package user

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"github.com/1024casts/snake/internal/cache/user"
	"github.com/1024casts/snake/internal/model"
)

// StatRepo 定义用户仓库接口
type StatRepo interface {
	IncrFollowCount(db *gorm.DB, userID uint64, step int) error
	IncrFollowerCount(db *gorm.DB, userID uint64, step int) error
	GetUserStatByID(db *gorm.DB, userID uint64) (*model.UserStatModel, error)
	GetUserStatByIDs(db *gorm.DB, userID []uint64) (map[uint64]*model.UserStatModel, error)
}

// userRepo 用户仓库
type userStatRepo struct {
	userCache *user.Cache
}

// NewUserStatRepo 实例化用户仓库
func NewUserStatRepo() StatRepo {
	return &userStatRepo{
		userCache: user.NewUserCache(),
	}
}

// IncrFollowCount 增加关注数
func (repo *userStatRepo) IncrFollowCount(db *gorm.DB, userID uint64, step int) error {
	err := db.Exec("insert into user_stat set user_id=?, follow_count=1, created_at=? on duplicate key update "+
		"follow_count=follow_count+?, updated_at=?",
		userID, time.Now(), step, time.Now()).Error
	if err != nil {
		return errors.Wrap(err, "[user_stat_repo] incr user follow count")
	}
	return nil
}

// IncrFollowerCount 增加粉丝数
func (repo *userStatRepo) IncrFollowerCount(db *gorm.DB, userID uint64, step int) error {
	err := db.Exec("insert into user_stat set user_id=?, follower_count=1, created_at=? on duplicate key update "+
		"follower_count=follower_count+?, updated_at=?",
		userID, time.Now(), step, time.Now()).Error
	if err != nil {
		return errors.Wrap(err, "[user_stat_repo] incr user follower count")
	}
	return nil
}

// GetUserStatByID 获取用户统计数据
func (repo *userStatRepo) GetUserStatByID(db *gorm.DB, userID uint64) (*model.UserStatModel, error) {
	userStat := model.UserStatModel{}
	err := db.Where("user_id = ?", userID).First(&userStat).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.Wrap(err, "[user_stat_repo] get user stat err")
	}

	return &userStat, nil
}

// GetUserStatByIDs 批量获取用户统计数据
func (repo *userStatRepo) GetUserStatByIDs(db *gorm.DB, userID []uint64) (map[uint64]*model.UserStatModel, error) {
	userStats := make([]*model.UserStatModel, 0)
	retMap := make(map[uint64]*model.UserStatModel)

	err := db.Where("user_id in (?)", userID).Find(&userStats).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return retMap, errors.Wrap(err, "[user_stat_repo] get user stat err")
	}

	for _, v := range userStats {
		retMap[v.UserID] = v
	}

	return retMap, nil
}
