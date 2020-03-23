package model

import (
	"sync"
	"time"

	"github.com/1024casts/snake/pkg/auth"
	validator "github.com/go-playground/validator/v10"
)

// User represents a registered user.
type UserModel struct {
	BaseModel
	Username      string    `json:"username" gorm:"column:username;not null" binding:"required" validate:"min=1,max=32"`
	Password      string    `json:"password" gorm:"column:password;not null" binding:"required" validate:"min=5,max=128"`
	Phone         int       `gorm:"column:phone" json:"phone"`
	Avatar        string    `gorm:"column:avatar" json:"avatar"`
	Sex           int       `gorm:"column:sex" json:"sex"`
	LastLoginIP   string    `gorm:"column:last_login_ip" json:"last_login_ip"`
	LastLoginTime time.Time `gorm:"column:last_login_time" json:"last_login_time"`
}

// Validate the fields.
func (u *UserModel) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

type UserInfo struct {
	Id        uint64 `json:"id" example:"1"`
	Username  string `json:"username" example:"张三"`
	Password  string `json:"password" example:"9dXd13#k$1123!kln"`
	CreatedAt string `json:"createdAt" example:"2020-03-23 20:00:00"`
	UpdatedAt string `json:"updatedAt" example:"2020-03-23 20:00:00"`
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
