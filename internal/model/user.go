package model

import (
	"sync"
	"time"

	"github.com/1024casts/snake/pkg/auth"

	validator "github.com/go-playground/validator/v10"
)

// UserModel User represents a registered user.
type UserModel struct {
	ID        uint64    `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"-"`
	Username  string    `json:"username" gorm:"column:username;not null" binding:"required" validate:"min=1,max=32"`
	Password  string    `json:"password" gorm:"column:password;not null" binding:"required" validate:"min=5,max=128"`
	Phone     int       `gorm:"column:phone" json:"phone"`
	Email     string    `gorm:"column:email" json:"email"`
	Avatar    string    `gorm:"column:avatar" json:"avatar"`
	Sex       int       `gorm:"column:sex" json:"sex"`
	CreatedAt time.Time `gorm:"column:created_at" json:"-"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"-"`
}

// Validate the fields.
func (u *UserModel) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

// UserFollow 关注
type UserFollow struct {
	FollowNum int `json:"follow_num"` // 关注数
	FansNum   int `json:"fans_num"`   // 粉丝数
	IsFollow  int `json:"is_follow"`  // 是否关注 1:是 0:否
	IsFans    int `json:"is_fans"`    // 是否是粉丝 1:是 0:否
}

// UserInfo 对外暴露的结构体
type UserInfo struct {
	ID         uint64      `json:"id" example:"1"`
	Username   string      `json:"username" example:"张三"`
	Avatar     string      `json:"avatar"`
	Sex        int         `json:"sex"`
	UserFollow *UserFollow `json:"user_follow"`
}

// TableName 表名
func (u *UserModel) TableName() string {
	return "users"
}

// UserList 用户列表结构体
type UserList struct {
	Lock  *sync.Mutex
	IDMap map[uint64]*UserInfo
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
