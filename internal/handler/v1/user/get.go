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

	// get the underlying request context
	ctx := c.Request.Context()

	// create a done channel to tell the request it's done
	doneChan := make(chan interface{})
	// create a err channel
	errChan := make(chan error)

	// here you put the actual work needed for the request
	// and then send the doneChan with the status and body
	// to finish the request by writing the response
	go func() {
		userID := cast.ToUint64(c.Param("id"))
		if userID == 0 {
			errChan <- errcode.ErrInvalidParam
			return
		}

		// Get the user by the `user_id` from the database.
		u, err := service.Svc.Users().GetUserByID(c.Request.Context(), userID)
		if errors.Is(err, repository.ErrNotFound) {
			log.Errorf("get user info err: %+v", err)
			errChan <- ecode.ErrUserNotFound
			return
		}
		if err != nil {
			errChan <- errcode.ErrInternalServer.WithDetails(err.Error())
			return
		}

		doneChan <- u
	}()

	// non-blocking select on two channels see if the request
	// times out or finishes
	select {
	// if the context is done it timed out or was canceled
	// so don't return anything
	case <-ctx.Done():
		return
	// if err is not nil return error response
	case err := <-errChan:
		response.Error(c, err)
	// if the request finished then finish the request by
	// writing the response
	case resp := <-doneChan:
		response.Success(c, resp)
	}
}
