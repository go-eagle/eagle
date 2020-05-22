package model

import "time"

// UserFollowModel 关注表
type UserFollowModel struct {
	ID          uint64    `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"-"`
	FollowedUID uint64    `gorm:"column:followed_uid" json:"followed_uid"`
	Status      int       `gorm:"column:status" json:"status"`
	UserID      uint64    `gorm:"column:user_id" json:"user_id"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"-"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"-"`
}

// TableName sets the insert table name for this struct type
func (u *UserFollowModel) TableName() string {
	return "user_follow"
}
