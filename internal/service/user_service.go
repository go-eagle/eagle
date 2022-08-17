package service

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/go-eagle/eagle/internal/repository"

	"github.com/pkg/errors"

	"github.com/go-eagle/eagle/internal/model"
	"github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/auth"
	"github.com/go-eagle/eagle/pkg/log"
)

// UserService define interface func
type UserService interface {
	Register(ctx context.Context, username, email, password string) error
	EmailLogin(ctx context.Context, email, password string) (tokenStr string, err error)
	PhoneLogin(ctx context.Context, phone int64, verifyCode int) (tokenStr string, err error)
	GetUserByID(ctx context.Context, id uint64) (*model.UserBaseModel, error)
	GetUserInfoByID(ctx context.Context, id uint64) (*model.UserInfo, error)
	GetUserByPhone(ctx context.Context, phone int64) (*model.UserBaseModel, error)
	GetUserByEmail(ctx context.Context, email string) (*model.UserBaseModel, error)
	UpdateUser(ctx context.Context, id uint64, userMap map[string]interface{}) error
	BatchGetUsers(ctx context.Context, userID uint64, userIDs []uint64) ([]*model.UserInfo, error)
}

type userService struct {
	repo repository.Repository
}

var _ UserService = (*userService)(nil)

func newUsers(svc *service) *userService {
	return &userService{repo: svc.repo}
}

// Register 注册用户
func (s *userService) Register(ctx context.Context, username, email, password string) error {
	pwd, err := auth.HashAndSalt(password)
	if err != nil {
		return errors.Wrapf(err, "encrypt password err")
	}

	u := model.UserBaseModel{
		Username:  username,
		Password:  pwd,
		Email:     email,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
	isExist, err := s.repo.UserIsExist(&u)
	if err != nil {
		return errors.Wrapf(err, "create user")
	}
	if isExist {
		return errors.New("用户已存在")
	}
	_, err = s.repo.CreateUser(ctx, &u)
	if err != nil {
		return errors.Wrapf(err, "create user")
	}
	return nil
}

// EmailLogin 邮箱登录
func (s *userService) EmailLogin(ctx context.Context, email, password string) (tokenStr string, err error) {
	u, err := s.GetUserByEmail(ctx, email)
	if err != nil {
		return "", errors.Wrapf(err, "get user info err by email")
	}

	// ComparePasswords the login password with the user password.
	if !auth.ComparePasswords(u.Password, password) {
		return "", errors.New("invalid password")
	}

	// 签发签名 Sign the json web token.
	payload := map[string]interface{}{"user_id": u.ID, "username": u.Username}
	tokenStr, err = app.Sign(ctx, payload, app.Conf.JwtSecret, 86400)
	if err != nil {
		return "", errors.Wrapf(err, "gen token sign err")
	}

	return tokenStr, nil
}

// PhoneLogin 邮箱登录
func (s *userService) PhoneLogin(ctx context.Context, phone int64, verifyCode int) (tokenStr string, err error) {
	// 如果是已经注册用户，则通过手机号获取用户信息
	u, err := s.GetUserByPhone(ctx, phone)
	if err != nil {
		return "", errors.Wrapf(err, "[login] get u info err")
	}

	// 否则新建用户信息, 并取得用户信息
	if u.ID == 0 {
		u := model.UserBaseModel{
			Phone:    phone,
			Username: strconv.Itoa(int(phone)),
		}
		u.ID, err = s.repo.CreateUser(ctx, &u)
		if err != nil {
			return "", errors.Wrapf(err, "[login] create user err")
		}
	}

	// 签发签名 Sign the json web token.
	payload := map[string]interface{}{"user_id": u.ID, "username": u.Username}
	tokenStr, err = app.Sign(ctx, payload, app.Conf.JwtSecret, 86400)
	if err != nil {
		return "", errors.Wrapf(err, "[login] gen token sign err")
	}
	return tokenStr, nil
}

// UpdateUser update user info
func (s *userService) UpdateUser(ctx context.Context, id uint64, userMap map[string]interface{}) error {
	err := s.repo.UpdateUser(ctx, id, userMap)

	if err != nil {
		return err
	}

	return nil
}

// GetUserByID 获取单条用户信息
func (s *userService) GetUserByID(ctx context.Context, id uint64) (*model.UserBaseModel, error) {
	return s.repo.GetUser(ctx, id)
}

// GetUserByID 获取单条用户信息
func (s *userService) GetUserStatByID(ctx context.Context, id uint64) (*model.UserStatModel, error) {
	return s.repo.GetUserStatByID(ctx, id)
}

// GetUserInfoByID 获取组装好的用户数据
func (s *userService) GetUserInfoByID(ctx context.Context, id uint64) (*model.UserInfo, error) {
	userInfos, err := s.BatchGetUsers(ctx, id, []uint64{id})
	if err != nil {
		return nil, err
	}
	return userInfos[0], nil
}

// BatchGetUsers 批量获取用户信息
// 1. 处理关注和被关注状态
// 2. 获取关注和粉丝数据
func (s *userService) BatchGetUsers(ctx context.Context, userID uint64, userIDs []uint64) ([]*model.UserInfo, error) {
	infos := make([]*model.UserInfo, 0)
	// 获取当前用户信息
	curUser, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "[user_service] get one user err")
	}

	// 批量获取用户信息
	users, err := s.repo.GetUsersByIds(ctx, userIDs)
	if err != nil {
		return nil, errors.Wrap(err, "[user_service] batch get user err")
	}

	wg := sync.WaitGroup{}

	errChan := make(chan error, 1)
	finished := make(chan bool, 1)

	// 获取自己对关注列表的关注状态
	userFollowMap, err := s.repo.GetFollowByUIds(ctx, userID, userIDs)
	if err != nil {
		errChan <- err
	}

	// 获取自己对关注列表的被关注状态
	userFansMap, err := s.repo.GetFansByUIds(ctx, userID, userIDs)
	if err != nil {
		errChan <- err
	}

	// 获取用户统计
	userStatMap, err := s.repo.GetUserStatByIDs(ctx, userIDs)
	if err != nil {
		errChan <- err
	}

	var m sync.Map

	// 并行处理
	for _, u := range users {
		wg.Add(1)
		go func(u *model.UserBaseModel) {
			defer wg.Done()

			isFollow := 0
			_, ok := userFollowMap[u.ID]
			if ok {
				isFollow = 1
			}

			isFollowed := 0
			_, ok = userFansMap[u.ID]
			if ok {
				isFollowed = 1
			}

			userStatMap, ok := userStatMap[u.ID]
			if !ok {
				userStatMap = nil
			}

			transInput := &TransferUserInput{
				CurUser:  curUser,
				User:     u,
				UserStat: userStatMap,
				IsFollow: isFollow,
				IsFans:   isFollowed,
			}
			userInfo := TransferUser(transInput)
			if err != nil {
				errChan <- err
				return
			}
			m.Store(u.ID, &userInfo)
		}(u)
	}

	go func() {
		wg.Wait()
		close(finished)
	}()

	select {
	case <-finished:
	case err := <-errChan:
		log.Warnf("[user_service] batch get user err chan: %v", err)
		return nil, err
	case <-time.After(3 * time.Second):
		return nil, fmt.Errorf("list users timeout after 3 seconds")
	}

	// 保证顺序
	for _, u := range users {
		info, _ := m.Load(u.ID)
		infos = append(infos, info.(*model.UserInfo))
	}

	return infos, nil
}

func (s *userService) GetUserByPhone(ctx context.Context, phone int64) (*model.UserBaseModel, error) {
	userModel, err := s.repo.GetUserByPhone(ctx, phone)
	if err != nil {
		return userModel, errors.Wrapf(err, "get user info err from db by phone: %d", phone)
	}

	return userModel, nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*model.UserBaseModel, error) {
	userModel, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return userModel, errors.Wrapf(err, "get user info err from db by email: %s", email)
	}

	return userModel, nil
}
