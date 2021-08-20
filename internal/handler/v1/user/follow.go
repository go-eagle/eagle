package user

import (
	"context"

	"github.com/go-eagle/eagle/internal/service"

	"github.com/go-eagle/eagle/internal/ecode"

	"github.com/gin-gonic/gin"

	"github.com/go-eagle/eagle/api"
	"github.com/go-eagle/eagle/pkg/errcode"
	"github.com/go-eagle/eagle/pkg/log"
)

// Follow 关注/取消关注
// @Summary 通过用户id关注/取消关注用户
// @Description Get an user by user id
// @Tags 用户
// @Accept  json
// @Produce  json
// @Param user_id body string true "用户id"
// @Success 200 {object} model.UserInfo "用户信息"
// @Router /users/follow [post]
func Follow(c *gin.Context) {
	var req FollowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warnf("follow bind param err: %v", err)
		response.Error(c, errcode.ErrBind.WithDetails(err.Error()))
		return
	}

	// Get the user by the `user_id` from the database.
	_, err := service.UserSvc.GetUserByID(c, req.UserID)
	if err != nil {
		response.Error(c, ecode.ErrUserNotFound.WithDetails(err.Error()))
		return
	}

	userID := api.GetUserID(c)
	// 不能关注自己
	if userID == req.UserID {
		response.Error(c, ecode.ErrUserNotFound)
		return
	}

	// 检查是否已经关注过
	isFollowed := service.UserSvc.IsFollowing(context.TODO(), userID, req.UserID)
	if isFollowed {
		response.Error(c, errcode.Success)
		return
	}

	if isFollowed {
		// 取消关注
		err = service.UserSvc.Unfollow(context.TODO(), userID, req.UserID)
		if err != nil {
			log.Warnf("[follow] cancel user follow err: %v", err)
			response.Error(c, errcode.ErrInternalServer.WithDetails(err.Error()))
			return
		}
	} else {
		// 添加关注
		err = service.UserSvc.Follow(context.TODO(), userID, req.UserID)
		if err != nil {
			log.Warnf("[follow] add user follow err: %v", err)
			response.Error(c, errcode.ErrInternalServer.WithDetails(err.Error()))
			return
		}
	}

	response.Success(c, nil)
}
