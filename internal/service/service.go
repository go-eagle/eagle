package service

import (
	"github.com/1024casts/snake/internal/service/user"
)

var (
	// Svc global service var
	Svc *Service
)

// Service struct
type Service struct {
	userSvc user.UserService
}

// New init service
func New() (s *Service) {
	s = &Service{
		userSvc: user.NewUserService(),
	}
	return s
}

// UserSvc return user service
func (s *Service) UserSvc() user.UserService {
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
