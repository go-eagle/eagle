package service

import (
	"github.com/1024casts/snake/model"
	"github.com/1024casts/snake/repository"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// 直接初始化，可以避免在使用时再实例化
var UserService = NewUserService()

// 校验码服务，生成校验码和获得校验码
type userService struct {
	userRepo repository.UserRepo
}

func NewUserService() *userService {
	return &userService{
		userRepo: repository.NewUserRepo(),
	}
}

func (srv *userService) CreateUser(user model.UserModel) (id uint64, err error) {
	id, err = srv.userRepo.CreateUser(model.GetDB(), user)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (srv *userService) UpdateUser(userMap map[string]interface{}, id uint64) error {
	err := srv.userRepo.Update(userMap, id)

	if err != nil {
		return err
	}

	return nil
}

func (srv *userService) GetUserById(id uint64) (*model.UserModel, error) {
	userModel, err := srv.userRepo.GetUserById(id)
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
		retMap[v.Id] = v
	}

	return retMap, nil
}

func (srv *userService) GetUserByPhone(phone int) (*model.UserModel, error) {
	userModel, err := srv.userRepo.GetUserByPhone(phone)
	if err != nil || gorm.IsRecordNotFoundError(err) {
		return userModel, errors.Wrapf(err, "get user info err from db by phone: %s", phone)
	}

	return userModel, nil
}
