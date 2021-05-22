package service

import (
	"github.com/1024casts/snake/internal/dao"
	"github.com/1024casts/snake/internal/model"
	"github.com/1024casts/snake/pkg/conf"
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
	c   *conf.Config
	dao *dao.Dao
}

// New init service
func New(c *conf.Config) (s *Service) {
	db := model.GetDB()
	s = &Service{
		c:   c,
		dao: dao.New(db),
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
	s.dao.Close()
}
