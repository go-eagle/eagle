package model

import (
	"sync"

	"github.com/1024casts/snake/pkg/auth"
	"github.com/1024casts/snake/pkg/errno"
	"github.com/1024casts/snake/pkg/valid"
	"github.com/pkg/errors"
)

// User represents a registered user.
type UserModel struct {
	BaseModel
	Username string `json:"username" gorm:"column:username;not null" binding:"required" validate:"min=1,max=32"`
	Password string `json:"password" gorm:"column:password;not null" binding:"required" validate:"min=5,max=128"`
	DefaultValidateChecker
}

func NewUser(request UserRequest) (*UserModel, error) {
	user := &UserModel{}

	if err := user.UpdateFromRequest(request); err != nil {
		return nil, errors.WithStack(err)
	}

	if err := user.Validate(); err != nil {
		return nil, errors.WithStack(err)
	}

	return user, nil
}

// Validate the fields.
//func (u *UserModel) Validate() error {
//	validate := validator.New()
//	return validate.Struct(u)
//}

func (user *UserModel) Validate() error {
	if valid.IsZero(
		user.Username,
	) {
		return errors.Wrap(errno.ErrParam, "webService")
	}
	user.SetValidated()
	return nil
}

type UserRequest struct {
	Username string
}

func (user *UserModel) UpdateFromRequest(request UserRequest) error {

	user.Username = request.Username

	return nil
}

type UserInfo struct {
	Id        uint64 `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func (c *UserModel) TableName() string {
	return "tb_users"
}

type UserList struct {
	Lock  *sync.Mutex
	IdMap map[uint64]*UserInfo
}

// Token represents a JSON web token.
type Token struct {
	Token string `json:"token"`
}

// Compare with the plain text password. Returns true if it's the same as the encrypted one (in the `User` struct).
func (u *UserModel) Compare(pwd string) (err error) {
	err = auth.Compare(u.Password, pwd)
	return
}

// Encrypt the user password.
func (u *UserModel) Encrypt() (err error) {
	u.Password, err = auth.Encrypt(u.Password)
	return
}
