package user

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"

	"github.com/go-eagle/eagle/internal/ecode"
	"github.com/go-eagle/eagle/internal/repository"
	"github.com/go-eagle/eagle/internal/service"
	"github.com/go-eagle/eagle/pkg/errcode"
	"github.com/go-eagle/eagle/pkg/log"
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

	userID := cast.ToUint64(c.Param("id"))
	if userID == 0 {
		response.Error(c, errcode.ErrInvalidParam)
		return
	}

	// Get the user by the `user_id` from the database.
	u, err := service.Svc.Users().GetUserByID(c.Request.Context(), userID)
	if errors.Is(err, repository.ErrNotFound) {
		log.Errorf("get user info err: %+v", err)
		response.Error(c, ecode.ErrUserNotFound)
		return
	}
	if err != nil {
		response.Error(c, errcode.ErrInternalServer.WithDetails(err.Error()))
		return
	}

	//time.Sleep(5 * time.Second)

	response.Success(c, u)
}
