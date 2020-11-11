package service

import (
	"github.com/1024casts/snake/internal/service/user"
	"github.com/1024casts/snake/pkg/conf"
)

var (
	// Svc global service var
	Svc *Service
)

// Service struct
type Service struct {
	c       *conf.Config
	userSvc user.IUserService
}

// New init service
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:       c,
		userSvc: user.NewUserService(),
	}
	return s
}

// UserSvc return user service
func (s *Service) UserSvc() user.IUserService {
	return s.userSvc
}

// Ping service
func (s *Service) Ping() error {
	return nil
}

// Close service
func (s *Service) Close() {
	s.userSvc.Close()
}
