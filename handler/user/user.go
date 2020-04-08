package user

import (
	"github.com/1024casts/snake/model"
)

// CreateRequest 创建用户请求
type CreateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// CreateResponse 创建用户响应
type CreateResponse struct {
	Username string `json:"username"`
}

// PhoneLoginCredentials 手机登录
type PhoneLoginCredentials struct {
	Phone      int `json:"phone" form:"phone" binding:"required" example:"13010002000"`
	VerifyCode int `json:"verify_code" form:"verify_code" binding:"required" example:"120110"`
}

// UpdateRequest 更新请求
type UpdateRequest struct {
	Avatar string `json:"avatar"`
	Sex    int    `json:"sex"`
}

// SwaggerListResponse 文档
type SwaggerListResponse struct {
	TotalCount uint64           `json:"totalCount"`
	UserList   []model.UserInfo `json:"userList"`
}
