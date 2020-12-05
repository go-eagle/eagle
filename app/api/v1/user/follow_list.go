package user

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/1024casts/snake/app/api"
	"github.com/1024casts/snake/internal/service"
	"github.com/1024casts/snake/pkg/errno"
	"github.com/1024casts/snake/pkg/log"
)

// FollowList 关注列表
// @Summary 正在关注的用户列表
// @Description Get an user by user id
// @Tags 用户
// @Accept  json
// @Produce  json
// @Param user_id body string true "用户id"
// @Success 200 {object} model.UserInfo "用户信息"
// @Router /users/{id}/following [get]
func FollowList(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, _ := strconv.Atoi(userIDStr)

	curUserID := api.GetUserID(c)
	log.Infof("cur uid: %d", curUserID)

	_, err := service.Svc.UserSvc().GetUserByID(context.TODO(), uint64(userID))
	if err != nil {
		api.SendResponse(c, errno.ErrUserNotFound, nil)
		return
	}

	lastIDStr := c.DefaultQuery("last_id", "0")
	lastID, _ := strconv.Atoi(lastIDStr)
	limit := 10

	userFollowList, err := service.Svc.RelationSvc().GetFollowingUserList(context.TODO(), uint64(userID), uint64(lastID), limit+1)
	if err != nil {
		log.Warnf("get following user list err: %+v", err)
		api.SendResponse(c, errno.InternalServerError, nil)
		return
	}

	hasMore := 0
	pageValue := lastID
	if len(userFollowList) > limit {
		hasMore = 1
		userFollowList = userFollowList[0 : len(userFollowList)-1]
		pageValue = lastID + 1
	}

	var userIDs []uint64
	for _, v := range userFollowList {
		userIDs = append(userIDs, v.FollowedUID)
	}

	userOutList, err := service.Svc.UserSvc().BatchGetUsers(context.TODO(), curUserID, userIDs)
	if err != nil {
		log.Warnf("batch get users err: %v", err)
		api.SendResponse(c, errno.InternalServerError, nil)
		return
	}

	api.SendResponse(c, errno.OK, ListResponse{
		TotalCount: 0,
		HasMore:    hasMore,
		PageKey:    "last_id",
		PageValue:  pageValue,
		Items:      userOutList,
	})
}
