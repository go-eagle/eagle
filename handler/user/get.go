package user

import (
	"strconv"

	"github.com/1024casts/snake/pkg/errno"

	"github.com/pkg/errors"

	. "github.com/1024casts/snake/handler"
	"github.com/1024casts/snake/model"
	"github.com/1024casts/snake/service"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

// @Summary Get an user by the user identifier
// @Description Get an user by username
// @Tags user
// @Accept  json
// @Produce  json
// @Param username path string true "Username"
// @Success 200 {object} model.UserModel "{"code":0,"message":"OK","data":{"username":"kong","password":"$2a$10$E0kwtmtLZbwW/bDQ8qI8e.eHPqhQOW9tvjwpyo/p05f/f4Qvr3OmS"}}"
// @Router /user/{username} [get]
func Get(c *gin.Context) {
	log.Info("Get function called.")

	userIdStr := c.Param("user_id")
	userId, _ := strconv.Atoi(userIdStr)

	userModel := &model.UserModel{}
	userModel.Id = uint64(userId)
	userSrv := service.NewUserService()
	err := userSrv.GetUserById(userModel)
	if err != nil {
		log.Warnf("get user err, %v", errors.WithStack(err))
		SendResponse(c, errno.InternalServerError, nil)
		return
	}

	SendResponse(c, nil, userModel)
}
