package user

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"github.com/1024casts/snake/model"
	"github.com/1024casts/snake/pkg/auth"
	"github.com/1024casts/snake/pkg/token"
	"github.com/1024casts/snake/repository/user"
)

// Service 用户服务接口定义
// 使用大写的service对外保留方法
type Service interface {
	Register(ctx *gin.Context, username, email, password string) error
	EmailLogin(ctx *gin.Context, email, password string) (tokenStr string, err error)
	PhoneLogin(ctx *gin.Context, phone int, verifyCode int) (tokenStr string, err error)
	GetUserByID(id uint64) (*model.UserModel, error)
	GetUserListByIds(id []uint64) (map[uint64]*model.UserModel, error)
	GetUserByPhone(phone int) (*model.UserModel, error)
	GetUserByEmail(email string) (*model.UserModel, error)
	UpdateUser(id uint64, userMap map[string]interface{}) error
}

// UserSvc 直接初始化，可以避免在使用时再实例化
var UserSvc = NewUserService()

// 用小写的 service 实现接口中定义的方法
type userService struct {
	userRepo user.Repo
}

// NewUserService 实例化一个userService
// 通过 NewService 函数初始化 Service 接口
// 依赖接口，不要依赖实现，面向接口编程
func NewUserService() Service {
	return &userService{
		userRepo: user.NewUserRepo(),
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
	if err != nil {
		return userModel, errors.Wrapf(err, "get user info err from db by id: %d", id)
	}

	return userModel, nil
}

// 批量获取
func (srv *userService) GetUserListByIds(id []uint64) (map[uint64]*model.UserModel, error) {
	userModels, err := srv.userRepo.GetUsersByIds(id)
	retMap := make(map[uint64]*model.UserModel)

	if err != nil {
		return retMap, errors.Wrapf(err, "get user model err from db by id: %v", id)
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
