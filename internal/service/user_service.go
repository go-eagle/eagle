package service

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/pkg/errors"

	v1 "github.com/1024casts/snake/api/grpc/user/v1"
	"github.com/1024casts/snake/internal/conf"
	"github.com/1024casts/snake/internal/model"
	"github.com/1024casts/snake/pkg/app"
	"github.com/1024casts/snake/pkg/auth"
	"github.com/1024casts/snake/pkg/log"
)

// Register 注册用户
func (s *Service) Register(ctx context.Context, username, email, password string) error {
	pwd, err := auth.Encrypt(password)
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
	_, err = s.userDao.CreateUser(ctx, u)
	if err != nil {
		return errors.Wrapf(err, "create user")
	}
	return nil
}

// EmailLogin 邮箱登录
func (s *Service) EmailLogin(ctx context.Context, email, password string) (tokenStr string, err error) {
	u, err := s.GetUserByEmail(ctx, email)
	if err != nil {
		return "", errors.Wrapf(err, "get user info err by email")
	}

	// Compare the login password with the user password.
	err = auth.Compare(u.Password, password)
	if err != nil {
		return "", errors.Wrapf(err, "password compare err")
	}

	// 签发签名 Sign the json web token.
	payload := map[string]interface{}{"user_id": u.ID, "username": u.Username}
	tokenStr, err = app.Sign(ctx, payload, conf.Conf.App.JwtSecret, 86400)
	if err != nil {
		return "", errors.Wrapf(err, "gen token sign err")
	}

	return tokenStr, nil
}

// LoginByPhone phone login, grpc wrapper
func (s *Service) LoginByPhone(ctx context.Context, req *v1.PhoneLoginRequest) (reply *v1.PhoneLoginReply, err error) {
	tokenStr, err := s.PhoneLogin(ctx, req.Phone, int(req.VerifyCode))
	if err != nil {
		log.Warnf("[service.user] phone login err: %v, params: %v", err, req)
	}
	reply = &v1.PhoneLoginReply{
		Ret: tokenStr,
		Err: "",
	}
	return
}

// PhoneLogin 邮箱登录
func (s *Service) PhoneLogin(ctx context.Context, phone int64, verifyCode int) (tokenStr string, err error) {
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
		u.ID, err = s.userDao.CreateUser(ctx, u)
		if err != nil {
			return "", errors.Wrapf(err, "[login] create user err")
		}
	}

	// 签发签名 Sign the json web token.
	payload := map[string]interface{}{"user_id": u.ID, "username": u.Username}
	tokenStr, err = app.Sign(ctx, payload, conf.Conf.App.JwtSecret, 86400)
	if err != nil {
		return "", errors.Wrapf(err, "[login] gen token sign err")
	}
	return tokenStr, nil
}

// UpdateUser update user info
func (s *Service) UpdateUser(ctx context.Context, id uint64, userMap map[string]interface{}) error {
	err := s.userDao.UpdateUser(ctx, id, userMap)

	if err != nil {
		return err
	}

	return nil
}

// GetUserByID 获取单条用户信息
func (s *Service) GetUserByID(ctx context.Context, id uint64) (*model.UserBaseModel, error) {
	userModel, err := s.userDao.GetOneUser(ctx, id)
	if err != nil {
		return userModel, errors.Wrap(err, "")
	}

	return userModel, nil
}

// GetUserInfoByID 获取组装好的用户数据
func (s *Service) GetUserInfoByID(ctx context.Context, id uint64) (*model.UserInfo, error) {
	userInfos, err := s.BatchGetUsers(ctx, id, []uint64{id})
	if err != nil {
		return nil, err
	}
	return userInfos[0], nil
}

// BatchGetUsers 批量获取用户信息
// 1. 处理关注和被关注状态
// 2. 获取关注和粉丝数据
func (s *Service) BatchGetUsers(ctx context.Context, userID uint64, userIDs []uint64) ([]*model.UserInfo, error) {
	infos := make([]*model.UserInfo, 0)
	// 批量获取用户信息
	users, err := s.userDao.GetUsersByIds(ctx, userIDs)
	if err != nil {
		return nil, errors.Wrap(err, "[user_service] batch get user err")
	}

	// 获取当前用户信息
	curUser, err := s.userDao.GetOneUser(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "[user_service] get one user err")
	}

	// 保持原有id顺序
	ids := userIDs

	wg := sync.WaitGroup{}
	userList := model.UserList{
		Lock:  new(sync.Mutex),
		IDMap: make(map[uint64]*model.UserInfo, len(users)),
	}

	errChan := make(chan error, 1)
	finished := make(chan bool, 1)

	// 获取自己对关注列表的关注状态
	userFollowMap, err := s.userFollowDao.GetFollowByUIds(ctx, userID, userIDs)
	if err != nil {
		errChan <- err
	}

	// 获取自己对关注列表的被关注状态
	userFansMap, err := s.userFollowDao.GetFansByUIds(ctx, userID, userIDs)
	if err != nil {
		errChan <- err
	}

	// 获取用户统计
	userStatMap, err := s.userStatDao.GetUserStatByIDs(ctx, userIDs)
	if err != nil {
		errChan <- err
	}

	// 并行处理
	for _, u := range users {
		wg.Add(1)
		go func(u *model.UserBaseModel) {
			defer wg.Done()

			userList.Lock.Lock()
			defer userList.Lock.Unlock()

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
			userList.IDMap[u.ID] = userInfo
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
	}

	// 根据原有id合并数据
	for _, id := range ids {
		infos = append(infos, userList.IDMap[id])
	}

	return infos, nil
}

func (s *Service) GetUserByPhone(ctx context.Context, phone int64) (*model.UserBaseModel, error) {
	userModel, err := s.userDao.GetUserByPhone(ctx, phone)
	if err != nil {
		return userModel, errors.Wrapf(err, "get user info err from db by phone: %d", phone)
	}

	return userModel, nil
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*model.UserBaseModel, error) {
	userModel, err := s.userDao.GetUserByEmail(ctx, email)
	if err != nil {
		return userModel, errors.Wrapf(err, "get user info err from db by email: %s", email)
	}

	return userModel, nil
}
