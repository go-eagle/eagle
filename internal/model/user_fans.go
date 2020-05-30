package model

import "time"

// UserFansModel 粉丝表
type UserFansModel struct {
	ID          uint64    `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"-"`
	FollowerUID uint64    `gorm:"column:follower_uid" json:"follower_uid"`
	Status      int       `gorm:"column:status" json:"status"`
	UserID      uint64    `gorm:"column:user_id" json:"user_id"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"-"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"-"`
}

// TableName sets the insert table name for this struct type
func (u *UserFansModel) TableName() string {
	return "user_fans"
}
