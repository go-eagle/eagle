package relation

import (
	"context"

	"github.com/opentracing/opentracing-go"

	"github.com/pkg/errors"

	"github.com/1024casts/snake/internal/model"
	"github.com/1024casts/snake/internal/repository/user"
	"github.com/1024casts/snake/pkg/conf"
	"github.com/1024casts/snake/pkg/constvar"
	"github.com/1024casts/snake/pkg/log"
)

const (
	// FollowStatusNormal 关注状态-正常
	FollowStatusNormal int = 1 // 正常
	// FollowStatusDelete 关注状态-删除
	FollowStatusDelete = 0 // 删除
)

// 用于触发编译期的接口的合理性检查机制
// 如果relationService没有实现IRelationService,会在编译期报错
var _ IRelationService = (*relationService)(nil)

// IRelationService 关系服务接口定义
type IRelationService interface {
	Follow(ctx context.Context, userID uint64, followedUID uint64) error
	Unfollow(ctx context.Context, userID uint64, followedUID uint64) error
	IsFollowing(ctx context.Context, userID uint64, followedUID uint64) bool
	GetFollowingUserList(ctx context.Context, userID uint64, lastID uint64, limit int) ([]*model.UserFollowModel, error)
	GetFollowerUserList(ctx context.Context, userID uint64, lastID uint64, limit int) ([]*model.UserFansModel, error)

	Close()
}

// relationService 用小写的 service 实现接口中定义的方法
type relationService struct {
	c              *conf.Config
	tracer         opentracing.Tracer
	userRepo       user.BaseRepo
	userFollowRepo user.FollowRepo
	userStatRepo   user.StatRepo
}

// NewRelationService 实例化一个userService
// 通过 NewService 函数初始化 Service 接口
// 依赖接口，不要依赖实现，面向接口编程
func NewRelationService(c *conf.Config, tracer opentracing.Tracer) IRelationService {
	db := model.GetDB()
	return &relationService{
		c:              c,
		tracer:         tracer,
		userRepo:       user.NewUserRepo(db),
		userFollowRepo: user.NewUserFollowRepo(db),
		userStatRepo:   user.NewUserStatRepo(db, tracer),
	}
}

// IsFollowing 是否正在关注某用户
func (srv *relationService) IsFollowing(ctx context.Context, userID uint64, followedUID uint64) bool {
	userFollowModel := &model.UserFollowModel{}
	result := model.GetDB().
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
func (srv *relationService) Follow(ctx context.Context, userID uint64, followedUID uint64) error {
	db := model.GetDB()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 添加到关注表
	err := srv.userFollowRepo.CreateUserFollow(ctx, tx, userID, followedUID)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "insert into user follow err")
	}

	// 添加到粉丝表
	err = srv.userFollowRepo.CreateUserFans(ctx, tx, followedUID, userID)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "insert into user fans err")
	}

	// 添加关注数
	err = srv.userStatRepo.IncrFollowCount(ctx, tx, userID, 1)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "update user follow count err")
	}

	// 添加粉丝数
	err = srv.userStatRepo.IncrFollowerCount(ctx, tx, followedUID, 1)
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
func (srv *relationService) Unfollow(ctx context.Context, userID uint64, followedUID uint64) error {
	db := model.GetDB()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除关注
	err := srv.userFollowRepo.UpdateUserFollowStatus(ctx, tx, userID, followedUID, FollowStatusDelete)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "update user follow err")
	}

	// 删除粉丝
	err = srv.userFollowRepo.UpdateUserFansStatus(ctx, tx, followedUID, userID, FollowStatusDelete)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "update user follow err")
	}

	// 减少关注数
	err = srv.userStatRepo.IncrFollowCount(ctx, tx, userID, -1)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "update user follow count err")
	}

	// 减少粉丝数
	err = srv.userStatRepo.IncrFollowerCount(ctx, tx, followedUID, -1)
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
func (srv *relationService) GetFollowingUserList(ctx context.Context, userID uint64, lastID uint64, limit int) ([]*model.UserFollowModel, error) {
	if lastID == 0 {
		lastID = constvar.MaxID
	}
	userFollowList, err := srv.userFollowRepo.GetFollowingUserList(ctx, userID, lastID, limit)
	if err != nil {
		return nil, err
	}

	return userFollowList, nil
}

// GetFollowerUserList 获取粉丝用户列表
func (srv *relationService) GetFollowerUserList(ctx context.Context, userID uint64, lastID uint64, limit int) ([]*model.UserFansModel, error) {
	if lastID == 0 {
		lastID = constvar.MaxID
	}
	userFollowerList, err := srv.userFollowRepo.GetFollowerUserList(ctx, userID, lastID, limit)
	if err != nil {
		return nil, err
	}

	return userFollowerList, nil
}

// Close close all user repo
func (srv *relationService) Close() {
	srv.userRepo.Close()
	srv.userFollowRepo.Close()
	srv.userStatRepo.Close()
}
