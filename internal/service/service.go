package service

import (
	"github.com/go-eagle/eagle/internal/repository"
)

// Svc global var
var Svc Service

const (
	// DefaultLimit 默认分页数
	DefaultLimit = 50

	// MaxID 最大id
	MaxID = 0xffffffffffff

	// DefaultAvatar 默认头像 key
	DefaultAvatar = "default_avatar.png"
)

// Service define all service
type Service interface {
	Users() UserService
	Relations() RelationService
	SMS() SMSService
	VCode() VCodeService
}

// service struct
type service struct {
	repo repository.Repository
}

// New init service
func New(repo repository.Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Users() UserService {
	return newUsers(s)
}

func (s *service) Relations() RelationService {
	return newRelations(s)
}

func (s *service) SMS() SMSService {
	return newSMS(s)
}

func (s *service) VCode() VCodeService {
	return newVCode(s)
}
