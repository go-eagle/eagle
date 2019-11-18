package service

import (
	"github.com/1024casts/snake/pkg/db"

	"github.com/1024casts/snake/repository"

	"github.com/realsangil/apimonitor/pkg/amerr"

	"github.com/1024casts/snake/model"
)

// 放到哪里合适呢？
type userRequest struct {
	Page int `json:"page"`
}

type UserService interface {
	CreateUser(request model.UserModel) (*model.UserModel, *amerr.ErrorWithLanguage)
	GetUserById(data *model.UserModel) error
	DeleteUserById(webService *model.UserModel) *amerr.ErrorWithLanguage
	UpdateUserById(webService *model.UserModel, request model.UserModel) *amerr.ErrorWithLanguage
	GetUserList(request userRequest) ([]*model.UserModel, int, error)
}

type UserServiceImpl struct {
	userRepo repository.UserRepository
}

func NewUserService() UserService {
	return &UserServiceImpl{
		userRepo: repository.NewUserRepository(),
	}
}

func (srv *UserServiceImpl) CreateUser(request model.UserModel) (*model.UserModel, *amerr.ErrorWithLanguage) {
	return nil, amerr.GetErrorsFromCode(amerr.ErrInternalServer)
}

func (srv *UserServiceImpl) GetUserById(userModel *model.UserModel) error {

	//if valid.IsZero(userModel) {
	//	return amerr.GetErrorsFromCode(amerr.ErrInternalServer)
	//}

	if err := srv.userRepo.GetById(db.GetConnection(), userModel); err != nil {
		switch err {
		case db.ErrRecordNotFound:
			return err
		}
		return err
	}
	return nil
}

func (service *UserServiceImpl) DeleteUserById(request *model.UserModel) *amerr.ErrorWithLanguage {
	return nil
}

func (service *UserServiceImpl) UpdateUserById(user *model.UserModel, request model.UserModel) *amerr.ErrorWithLanguage {
	return nil
}

func (service *UserServiceImpl) GetUserList(request userRequest) ([]*model.UserModel, int, error) {

	items := make([]*model.UserModel, 0)
	totalCount, err := service.userRepo.List(db.GetConnection(), &items, db.ListFilter{
		Page:       request.Page,
		Conditions: map[string]interface{}{},
	}, db.Orders{
		{
			Field: "id",
			IsASC: false,
		},
	})

	if err != nil {
		return items, 0, err
	}
	return items, totalCount, nil
}

//func ListUser(username string, offset, limit int) ([]*model.UserInfo, uint64, error) {
//	infos := make([]*model.UserInfo, 0)
//	users, count, err := model.ListUser(username, offset, limit)
//	if err != nil {
//		return nil, count, err
//	}
//
//	var ids []uint64
//	for _, user := range users {
//		ids = append(ids, user.Id)
//	}
//
//	wg := sync.WaitGroup{}
//	userList := model.UserList{
//		Lock:  new(sync.Mutex),
//		IdMap: make(map[uint64]*model.UserInfo, len(users)),
//	}
//
//	errChan := make(chan error, 1)
//	finished := make(chan bool, 1)
//
//	// Improve query efficiency in parallel
//	for _, u := range users {
//		wg.Add(1)
//		go func(u *model.UserModel) {
//			defer wg.Done()
//
//			shortId, err := util.GenShortId()
//			if err != nil {
//				errChan <- err
//				return
//			}
//
//			userList.Lock.Lock()
//			defer userList.Lock.Unlock()
//			userList.IdMap[u.Id] = &model.UserInfo{
//				Id:        u.Id,
//				Username:  u.Username,
//				SayHello:  fmt.Sprintf("Hello %s", shortId),
//				Password:  u.Password,
//				CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
//				UpdatedAt: u.UpdatedAt.Format("2006-01-02 15:04:05"),
//			}
//		}(u)
//	}
//
//	go func() {
//		wg.Wait()
//		close(finished)
//	}()
//
//	select {
//	case <-finished:
//	case err := <-errChan:
//		return nil, count, err
//	}
//
//	for _, id := range ids {
//		infos = append(infos, userList.IdMap[id])
//	}
//
//	return infos, count, nil
//}
