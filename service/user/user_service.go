package user

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"github.com/1024casts/snake/model"
	"github.com/1024casts/snake/pkg/auth"
	"github.com/1024casts/snake/pkg/log"
	"github.com/1024casts/snake/pkg/token"
	"github.com/1024casts/snake/repository/user"
)

const (
	// FollowStatusNormal 关注状态-正常
	FollowStatusNormal int = 1 // 正常
	// FollowStatusDelete 关注状态-删除
	FollowStatusDelete = 0 // 删除

	// MaxID 最大id
	MaxID = 0xffffffffffff
)

// Service 用户服务接口定义
// 使用大写的service对外保留方法
type Service interface {
	Register(ctx *gin.Context, username, email, password string) error
	EmailLogin(ctx *gin.Context, email, password string) (tokenStr string, err error)
	PhoneLogin(ctx *gin.Context, phone int, verifyCode int) (tokenStr string, err error)
	GetUserByID(id uint64) (*model.UserModel, error)
	BatchGetUserListByIds(userID []uint64) (map[uint64]*model.UserModel, error)
	GetUserByPhone(phone int) (*model.UserModel, error)
	GetUserByEmail(email string) (*model.UserModel, error)
	UpdateUser(id uint64, userMap map[string]interface{}) error

	// 关注
	IsFollowedUser(userID uint64, followedUID uint64) bool
	AddUserFollow(userID uint64, followedUID uint64) error
	CancelUserFollow(userID uint64, followedUID uint64) error
	GetFollowingUserList(userID uint64, lastID uint64, limit int) ([]*model.UserFollowModel, error)
	GetFollowerUserList(userID uint64, lastID uint64, limit int) ([]*model.UserFansModel, error)
}

// Svc 直接初始化，可以避免在使用时再实例化
var Svc = NewUserService()

// 用小写的 service 实现接口中定义的方法
type userService struct {
	userRepo       user.Repo
	userFollowRepo user.FollowRepo
}

// NewUserService 实例化一个userService
// 通过 NewService 函数初始化 Service 接口
// 依赖接口，不要依赖实现，面向接口编程
func NewUserService() Service {
	return &userService{
		userRepo:       user.NewUserRepo(),
		userFollowRepo: user.NewUserFollowRepo(),
	}
}

// Register 注册用户
func (srv *userService) Register(ctx *gin.Context, username, email, password string) error {
	pwd, err := auth.Encrypt(password)
	if err != nil {
		return errors.Wrapf(err, "encrypt password err")
	}

	u := model.UserModel{
		Username:  username,
		Password:  pwd,
		Email:     email,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
	_, err = srv.userRepo.Create(model.GetDB(), u)
	if err != nil {
		return errors.Wrapf(err, "create user")
	}
	return nil
}

// EmailLogin 邮箱登录
func (srv *userService) EmailLogin(ctx *gin.Context, email, password string) (tokenStr string, err error) {
	u, err := srv.GetUserByEmail(email)
	if err != nil {
		return "", errors.Wrapf(err, "get user info err by email")
	}

	// Compare the login password with the user password.
	err = auth.Compare(u.Password, password)
	if err != nil {
		return "", errors.Wrapf(err, "password compare err")
	}

	// 签发签名 Sign the json web token.
	tokenStr, err = token.Sign(ctx, token.Context{UserID: u.ID, Username: u.Username}, "")
	if err != nil {
		return "", errors.Wrapf(err, "gen token sign err")
	}

	return tokenStr, nil
}

// PhoneLogin 邮箱登录
func (srv *userService) PhoneLogin(ctx *gin.Context, phone int, verifyCode int) (tokenStr string, err error) {
	// 如果是已经注册用户，则通过手机号获取用户信息
	u, err := srv.GetUserByPhone(phone)
	if err != nil {
		return "", errors.Wrapf(err, "[login] get u info err")
	}

	// 否则新建用户信息, 并取得用户信息
	if u.ID == 0 {
		u := model.UserModel{
			Phone:    phone,
			Username: strconv.Itoa(phone),
		}
		u.ID, err = srv.userRepo.Create(model.GetDB(), u)
		if err != nil {
			return "", errors.Wrapf(err, "[login] create user err")
		}
	}

	// 签发签名 Sign the json web token.
	tokenStr, err = token.Sign(ctx, token.Context{UserID: u.ID, Username: u.Username}, "")
	if err != nil {
		return "", errors.Wrapf(err, "[login] gen token sign err")
	}
	return tokenStr, nil
}

func (srv *userService) UpdateUser(id uint64, userMap map[string]interface{}) error {
	err := srv.userRepo.Update(model.GetDB(), id, userMap)

	if err != nil {
		return err
	}

	return nil
}

func (srv *userService) GetUserByID(id uint64) (*model.UserModel, error) {
	userModel, err := srv.userRepo.GetUserByID(model.GetDB(), id)
	if err != nil && err != gorm.ErrRecordNotFound {
		return userModel, errors.Wrapf(err, "get user info err from db by id: %d", id)
	}

	return userModel, nil
}

// BatchGetUserListByIds 批量获取用户信息
func (srv *userService) BatchGetUserListByIds(userID []uint64) (map[uint64]*model.UserModel, error) {
	userModels, err := srv.userRepo.GetUsersByIds(userID)
	retMap := make(map[uint64]*model.UserModel)

	if err != nil {
		return retMap, errors.Wrapf(err, "get user model err from db by id: %v", userID)
	}

	for _, v := range userModels {
		retMap[v.ID] = v
	}

	return retMap, nil
}

func (srv *userService) GetUserByPhone(phone int) (*model.UserModel, error) {
	userModel, err := srv.userRepo.GetUserByPhone(model.GetDB(), phone)
	if err != nil || gorm.IsRecordNotFoundError(err) {
		return userModel, errors.Wrapf(err, "get user info err from db by phone: %d", phone)
	}

	return userModel, nil
}

func (srv *userService) GetUserByEmail(email string) (*model.UserModel, error) {
	userModel, err := srv.userRepo.GetUserByEmail(model.GetDB(), email)
	if err != nil || gorm.IsRecordNotFoundError(err) {
		return userModel, errors.Wrapf(err, "get user info err from db by email: %s", email)
	}

	return userModel, nil
}

// 获取用户关注
func (srv *userService) GetFollowUser(userID uint64, followedUID uint64) (*model.UserFollowModel, error) {
	userFollowModel := &model.UserFollowModel{}
	result := model.GetDB().
		Where("user_id=? AND followed_uid=? ", userID, followedUID).
		Find(userFollowModel)

	return userFollowModel, result.Error
}

// IsFollowedUser 是否关注过某用户
func (srv *userService) IsFollowedUser(userID uint64, followedUID uint64) bool {
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

// AddUserFollow 添加关注
func (srv *userService) AddUserFollow(userID uint64, followedUID uint64) error {
	db := model.GetDB()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 添加到关注表
	err := srv.userFollowRepo.CreateUserFollow(tx, userID, followedUID)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "insert into user follow err")
	}

	// 添加到粉丝表
	err = srv.userFollowRepo.CreateUserFans(tx, followedUID, userID)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "insert into user fans err")
	}

	// 添加关注数
	err = srv.userRepo.IncrFollowCount(tx, userID, 1)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "update user follow count err")
	}

	// 添加粉丝数
	err = srv.userRepo.IncrFollowerCount(tx, followedUID, 1)
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

// CancelUserFollow 取消用户关注
func (srv *userService) CancelUserFollow(userID uint64, followedUID uint64) error {
	db := model.GetDB()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除关注
	err := srv.userFollowRepo.UpdateUserFollowStatus(tx, userID, followedUID, FollowStatusDelete)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "update user follow err")
	}

	// 删除粉丝
	err = srv.userFollowRepo.UpdateUserFansStatus(tx, followedUID, userID, FollowStatusDelete)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "update user follow err")
	}

	// 减少关注数
	err = srv.userRepo.IncrFollowCount(tx, userID, -1)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "update user follow count err")
	}

	// 减少粉丝数
	err = srv.userRepo.IncrFollowerCount(tx, followedUID, -1)
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
func (srv *userService) GetFollowingUserList(userID uint64, lastID uint64, limit int) ([]*model.UserFollowModel, error) {
	if lastID == 0 {
		lastID = MaxID
	}
	userFollowList, err := srv.userFollowRepo.GetFollowingUserList(userID, lastID, limit)
	if err != nil {
		return nil, err
	}

	return userFollowList, nil
}

// GetFollowerUserList 获取粉丝用户列表
func (srv *userService) GetFollowerUserList(userID uint64, lastID uint64, limit int) ([]*model.UserFansModel, error) {
	if lastID == 0 {
		lastID = MaxID
	}
	userFollowerList, err := srv.userFollowRepo.GetFollowerUserList(userID, lastID, limit)
	if err != nil {
		return nil, err
	}

	return userFollowerList, nil
}
