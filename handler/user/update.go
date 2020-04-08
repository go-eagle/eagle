package user

import (
	"strconv"

	"github.com/1024casts/snake/service/user"

	. "github.com/1024casts/snake/handler"
	"github.com/1024casts/snake/pkg/errno"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

// Update 更新用户信息
// @Summary Update a user info by the user identifier
// @Description Update a user by ID
// @Tags 用户
// @Accept  json
// @Produce  json
// @Param id path uint64 true "The user's database id index num"
// @Param user body model.UserModel true "The user info"
// @Success 200 {object} handler.Response "{"code":0,"message":"OK","data":null}"
// @Router /users/{id} [put]
func Update(c *gin.Context) {
	log.Info("Update function called.")
	// Get the user id from the url parameter.
	userID, _ := strconv.Atoi(c.Param("id"))

	// Binding the user data.
	var req UpdateRequest
	if err := c.Bind(&req); err != nil {
		SendResponse(c, errno.ErrBind, nil)
		return
	}

	userMap := make(map[string]interface{})
	userMap["avatar"] = req.Avatar
	userMap["sex"] = req.Sex
	err := user.UserService.UpdateUser(userMap, uint64(userID))
	if err != nil {
		log.Warnf("[user] update user err, %v", err)
		SendResponse(c, errno.InternalServerError, nil)
		return
	}

	SendResponse(c, nil, userID)
}
