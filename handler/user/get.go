package user

import (
	"strconv"

	"github.com/1024casts/snake/service/user"

	. "github.com/1024casts/snake/handler"
	"github.com/1024casts/snake/pkg/errno"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

// Get 获取用户信息
// @Summary 通过用户id获取用户信息
// @Description Get an user by user id
// @Tags 用户
// @Accept  json
// @Produce  json
// @Param id path string true "用户id"
// @Success 200 {object} model.UserInfo "用户信息"
// @Router /users/:id [get]
func Get(c *gin.Context) {
	log.Info("Get function called.")

	userIDStr := c.Param("id")
	if userIDStr == "" {
		SendResponse(c, errno.ErrParam, nil)
		return
	}
	userID, _ := strconv.Atoi(userIDStr)

	// Get the user by the `user_id` from the database.
	u, err := user.UserService.GetUserByID(uint64(userID))
	if err != nil {
		SendResponse(c, errno.ErrUserNotFound, nil)
		return
	}

	SendResponse(c, nil, u)
}
