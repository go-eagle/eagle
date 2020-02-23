package user

import (
	. "github.com/1024casts/snake/handler"
	"github.com/1024casts/snake/pkg/errno"
	"github.com/1024casts/snake/service"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

// @Summary 获取自己的个人信息 Done
// @Description Get an user by user id
// @Tags 用户
// @Accept  json
// @Produce  json
// @Success 200 {object} model.UserInfo "个人用户信息"
// @Router /users/me [get]
func Me(c *gin.Context) {
	log.Info("Me function called.")

	userId := GetUserId(c)

	// Get the user by the `user_id` from the database.
	user, err := service.UserService.GetUserById(userId)
	if err != nil {
		SendResponse(c, errno.ErrUserNotFound, nil)
		return
	}

	SendResponse(c, nil, user)
}
