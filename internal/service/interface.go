package service

import (
	"context"

	v1 "github.com/go-eagle/eagle/api/grpc/user/v1"
	"github.com/go-eagle/eagle/internal/model"
)

// 用于触发编译期的接口的合理性检查机制
var _ service = (*Service)(nil)

// service 接口定义
type service interface {
	// user
	Register(ctx context.Context, username, email, password string) error
	EmailLogin(ctx context.Context, email, password string) (tokenStr string, err error)
	PhoneLogin(ctx context.Context, phone int64, verifyCode int) (tokenStr string, err error)
	LoginByPhone(ctx context.Context, req *v1.PhoneLoginRequest) (reply *v1.PhoneLoginReply, err error)
	GetUserByID(ctx context.Context, id uint64) (*model.UserBaseModel, error)
	GetUserInfoByID(ctx context.Context, id uint64) (*model.UserInfo, error)
	GetUserByPhone(ctx context.Context, phone int64) (*model.UserBaseModel, error)
	GetUserByEmail(ctx context.Context, email string) (*model.UserBaseModel, error)
	UpdateUser(ctx context.Context, id uint64, userMap map[string]interface{}) error
	BatchGetUsers(ctx context.Context, userID uint64, userIDs []uint64) ([]*model.UserInfo, error)

	// relation
	Follow(ctx context.Context, userID uint64, followedUID uint64) error
	Unfollow(ctx context.Context, userID uint64, followedUID uint64) error
	IsFollowing(ctx context.Context, userID uint64, followedUID uint64) bool
	GetFollowingUserList(ctx context.Context, userID uint64, lastID uint64, limit int) ([]*model.UserFollowModel, error)
	GetFollowerUserList(ctx context.Context, userID uint64, lastID uint64, limit int) ([]*model.UserFansModel, error)

	// sms
	SendSMS(phoneNumber string, verifyCode int) error

	// verify code
	GenLoginVCode(phone string) (int, error)
	CheckLoginVCode(phone int64, vCode int) bool
	GetLoginVCode(phone int64) (int, error)
}
