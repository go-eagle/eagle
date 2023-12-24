package service

import (
	"context"

	"github.com/pkg/errors"

	"github.com/go-eagle/eagle/internal/model"
	"github.com/go-eagle/eagle/internal/repository"
	"github.com/go-eagle/eagle/pkg/log"
)

const (
	// FollowStatusNormal 关注状态-正常
	FollowStatusNormal int = 1 // 正常
	// FollowStatusDelete 关注状态-删除
	FollowStatusDelete = 0 // 删除
)

// RelationService .
type RelationService interface {
	Follow(ctx context.Context, userID uint64, followedUID uint64) error
	Unfollow(ctx context.Context, userID uint64, followedUID uint64) error
	IsFollowing(ctx context.Context, userID uint64, followedUID uint64) bool
	GetFollowingUserList(ctx context.Context, userID uint64, lastID uint64, limit int) ([]*model.UserFollowModel, error)
	GetFollowerUserList(ctx context.Context, userID uint64, lastID uint64, limit int) ([]*model.UserFansModel, error)
}

type relationService struct {
	repo repository.Repository
}

var _ RelationService = (*relationService)(nil)

func newRelations(svc *service) *relationService {
	return &relationService{repo: svc.repo}
}

// IsFollowing 是否正在关注某用户
func (s *relationService) IsFollowing(ctx context.Context, userID uint64, followedUID uint64) bool {
	userFollowModel := &model.UserFollowModel{}
	db, _ := model.GetDB()
	result := db.
		Where("user_id=? AND followed_uid=? ", userID, followedUID).
		Find(userFollowModel)

	if err := result.Error; err != nil {
		log.Warnf("[user_service] get user follow err, %v", err)
		return false
	}

	if userFollowModel.ID > 0 && userFollowModel.Status == FollowStatusNormal {
		return true
	}

	return false
}

// Follow 关注目标用户
func (s *relationService) Follow(ctx context.Context, userID uint64, followedUID uint64) error {
	db, _ := model.GetDB()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 添加到关注表
	err := s.repo.CreateUserFollow(ctx, tx, userID, followedUID)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "insert into user follow err")
	}

	// 添加到粉丝表
	err = s.repo.CreateUserFans(ctx, tx, followedUID, userID)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "insert into user fans err")
	}

	// 添加关注数
	err = s.repo.IncrFollowCount(ctx, tx, userID, 1)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "update user follow count err")
	}

	// 添加粉丝数
	err = s.repo.IncrFollowerCount(ctx, tx, followedUID, 1)
	if err != nil {
		return errors.Wrap(err, "update user fans count err")
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "tx commit err")
	}

	return nil
}

// Unfollow 取消用户关注
func (s *relationService) Unfollow(ctx context.Context, userID uint64, followedUID uint64) error {
	db, _ := model.GetDB()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除关注
	err := s.repo.UpdateUserFollowStatus(ctx, tx, userID, followedUID, FollowStatusDelete)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "update user follow err")
	}

	// 删除粉丝
	err = s.repo.UpdateUserFansStatus(ctx, tx, followedUID, userID, FollowStatusDelete)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "update user follow err")
	}

	// 减少关注数
	err = s.repo.IncrFollowCount(ctx, tx, userID, -1)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "update user follow count err")
	}

	// 减少粉丝数
	err = s.repo.IncrFollowerCount(ctx, tx, followedUID, -1)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "update user fans count err")
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "tx commit err")
	}

	return nil
}

// GetFollowingUserList 获取正在关注的用户列表
func (s *relationService) GetFollowingUserList(ctx context.Context, userID uint64, lastID uint64, limit int) ([]*model.UserFollowModel, error) {
	if lastID == 0 {
		lastID = MaxID
	}
	userFollowList, err := s.repo.GetFollowingUserList(ctx, userID, lastID, limit)
	if err != nil {
		return nil, err
	}

	return userFollowList, nil
}

// GetFollowerUserList 获取粉丝用户列表
func (s *relationService) GetFollowerUserList(ctx context.Context, userID uint64, lastID uint64, limit int) ([]*model.UserFansModel, error) {
	if lastID == 0 {
		lastID = MaxID
	}
	userFollowerList, err := s.repo.GetFollowerUserList(ctx, userID, lastID, limit)
	if err != nil {
		return nil, err
	}

	return userFollowerList, nil
}
