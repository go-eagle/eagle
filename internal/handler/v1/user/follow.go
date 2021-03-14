package user

import (
	"context"

	"github.com/1024casts/snake/internal/service"

	"github.com/1024casts/snake/internal/ecode"

	"github.com/gin-gonic/gin"

	"github.com/1024casts/snake/api"
	"github.com/1024casts/snake/pkg/errno"
	"github.com/1024casts/snake/pkg/log"
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
		Response.Error(c, errno.ErrBind)
		return
	}

	// Get the user by the `user_id` from the database.
	_, err := service.UserSvc.GetUserByID(c, req.UserID)
	if err != nil {
		Response.Error(c, ecode.ErrUserNotFound)
		return
	}

	userID := api.GetUserID(c)
	// 不能关注自己
	if userID == req.UserID {
		Response.Error(c, ecode.ErrUserNotFound)
		return
	}

	// 检查是否已经关注过
	isFollowed := service.UserSvc.IsFollowing(context.TODO(), userID, req.UserID)
	if isFollowed {
		Response.Error(c, errno.Success)
		return
	}

	if isFollowed {
		// 取消关注
		err = service.UserSvc.Unfollow(context.TODO(), userID, req.UserID)
		if err != nil {
			log.Warnf("[follow] cancel user follow err: %v", err)
			Response.Error(c, errno.InternalServerError)
			return
		}
	} else {
		// 添加关注
		err = service.UserSvc.Follow(context.TODO(), userID, req.UserID)
		if err != nil {
			log.Warnf("[follow] add user follow err: %v", err)
			Response.Error(c, errno.InternalServerError)
			return
		}
	}

	Response.Success(c, nil)
}
