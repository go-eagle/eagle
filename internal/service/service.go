package service

import (
	"github.com/1024casts/snake/internal/conf"
	"github.com/1024casts/snake/internal/dao"
	"github.com/1024casts/snake/internal/model"
)

const (
	// DefaultLimit 默认分页数
	DefaultLimit = 50

	// MaxID 最大id
	MaxID = 0xffffffffffff

	// DefaultAvatar 默认头像 key
	DefaultAvatar = "default_avatar.png"
)

var (
	UserSvc  *Service
	VCodeSvc *Service
)

// Service struct
type Service struct {
	c             *conf.Config
	userDao       dao.BaseDao
	userFollowDao dao.FollowRepo
	userStatDao   dao.StatRepo
}

// New init service
func New(c *conf.Config) (s *Service) {
	db := model.GetDB()
	s = &Service{
		c:             c,
		userDao:       dao.NewUserDao(db),
		userFollowDao: dao.NewUserFollowRepo(db),
		userStatDao:   dao.NewUserStatRepo(db),
	}
	UserSvc = s
	VCodeSvc = s
	return s
}

// Ping service
func (s *Service) Ping() error {
	return nil
}

// Close service
func (s *Service) Close() {
	s.userDao.Close()
	s.userFollowDao.Close()
	s.userStatDao.Close()
}
