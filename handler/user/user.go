package user

import (
	"github.com/1024casts/snake/model"
)

type CreateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateResponse struct {
	Username string `json:"username"`
}

// 手机登录
type PhoneLoginCredentials struct {
	Phone      int `json:"phone" form:"phone" binding:"required" example:"13010002000"`
	VerifyCode int `json:"verify_code" form:"verify_code" binding:"required" example:"120110"`
}

type ListRequest struct {
	Username string `json:"username"`
	Offset   int    `json:"offset"`
	Limit    int    `json:"limit"`
}

type ListResponse struct {
	TotalCount uint64            `json:"totalCount"`
	UserList   []*model.UserInfo `json:"userList"`
}

type SwaggerListResponse struct {
	TotalCount uint64           `json:"totalCount"`
	UserList   []model.UserInfo `json:"userList"`
}
