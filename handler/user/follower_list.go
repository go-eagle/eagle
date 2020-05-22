package user

import (
	"strconv"

	"github.com/1024casts/snake/handler"
	"github.com/1024casts/snake/pkg/errno"
	"github.com/1024casts/snake/pkg/log"
	"github.com/1024casts/snake/service/user"
	"github.com/gin-gonic/gin"
)

// Get 粉丝列表
// @Summary 通过用户id关注用户
// @Description Get an user by user id
// @Tags 用户
// @Accept  json
// @Produce  json
// @Param user_id body string true "用户id"
// @Success 200 {object} model.UserInfo "用户信息"
// @Router /users/followers [get]
func FollowerList(c *gin.Context) {
	userID := handler.GetUserID(c)

	_, err := user.UserSvc.GetUserByID(userID)
	if err != nil {
		handler.SendResponse(c, errno.ErrUserNotFound, nil)
		return
	}

	lastIdStr := c.DefaultQuery("last_id", "0")
	lastID, _ := strconv.Atoi(lastIdStr)
	limit := 10

	userFollowList, err := user.UserSvc.GetFollowerUserList(userID, uint64(lastID), limit+1)
	if err != nil {
		log.Warnf("get follower user list err: %+v", err)
		handler.SendResponse(c, errno.InternalServerError, nil)
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

	userOutList, err := user.UserSvc.BatchGetUserListByIds(userIDs)
	if err != nil {
		handler.SendResponse(c, errno.InternalServerError, nil)
		return
	}

	handler.SendResponse(c, errno.OK, ListResponse{
		TotalCount: 0,
		HasMore:    hasMore,
		PageKey:    "last_id",
		PageValue:  pageValue,
		Items:      userOutList,
	})
	return
}
