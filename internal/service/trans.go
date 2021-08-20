package service

import (
	"github.com/go-eagle/eagle/internal/model"
)

// TransferUserInput 转换输入字段
type TransferUserInput struct {
	CurUser  *model.UserBaseModel
	User     *model.UserBaseModel
	UserStat *model.UserStatModel
	IsFollow int `json:"is_follow"`
	IsFans   int `json:"is_fans"`
}

// TransferUser 组装数据并输出
// 对外暴露的user结构，都应该经过此结构进行转换
func TransferUser(input *TransferUserInput) *model.UserInfo {
	if input.User == nil {
		return &model.UserInfo{}
	}

	return &model.UserInfo{
		ID:         input.User.ID,
		Username:   input.User.Username,
		Avatar:     input.User.Avatar, // todo: 转为url
		Sex:        input.User.Sex,
		UserFollow: transferUserFollow(input),
	}
}

// transferUserFollow 转换用户关注相关字段
func transferUserFollow(input *TransferUserInput) *model.UserFollow {
	followCount := 0
	if input.UserStat != nil {
		followCount = input.UserStat.FollowCount
	}
	followerCount := 0
	if input.UserStat != nil {
		followerCount = input.UserStat.FollowerCount
	}

	return &model.UserFollow{
		FollowNum: followCount,
		FansNum:   followerCount,
		IsFollow:  input.IsFollow,
		IsFans:    input.IsFans,
	}
}
