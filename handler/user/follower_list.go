package user

import (
	"strconv"

	"github.com/1024casts/snake/idl"
	"github.com/1024casts/snake/model"

	"github.com/gin-gonic/gin"

	"github.com/1024casts/snake/handler"
	"github.com/1024casts/snake/pkg/errno"
	"github.com/1024casts/snake/pkg/log"
	"github.com/1024casts/snake/service/user"
)

// FollowerList 粉丝列表
// @Summary 通过用户id关注用户
// @Description Get an user by user id
// @Tags 用户
// @Accept  json
// @Produce  json
// @Param user_id body string true "用户id"
// @Success 200 {object} model.UserInfo "用户信息"
// @Router /users/{id}/followers [get]
func FollowerList(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, _ := strconv.Atoi(userIDStr)

	curUser, err := user.Svc.GetUserByID(handler.GetUserID(c))
	if err != nil {
		handler.SendResponse(c, errno.ErrUserNotFound, nil)
		return
	}

	_, err = user.Svc.GetUserByID(uint64(userID))
	if err != nil {
		handler.SendResponse(c, errno.ErrUserNotFound, nil)
		return
	}

	lastIDStr := c.DefaultQuery("last_id", "0")
	lastID, _ := strconv.Atoi(lastIDStr)
	limit := 10

	userFollowerList, err := user.Svc.GetFollowerUserList(uint64(userID), uint64(lastID), limit+1)
	if err != nil {
		log.Warnf("get follower user list err: %+v", err)
		handler.SendResponse(c, errno.InternalServerError, nil)
		return
	}

	hasMore := 0
	pageValue := lastID
	if len(userFollowerList) > limit {
		hasMore = 1
		userFollowerList = userFollowerList[0 : len(userFollowerList)-1]
		pageValue = lastID + 1
	}

	var userIDs []uint64
	for _, v := range userFollowerList {
		userIDs = append(userIDs, v.FollowerUID)
	}

	userMap, err := user.Svc.BatchGetUserListByIds(userIDs)
	if err != nil {
		handler.SendResponse(c, errno.InternalServerError, nil)
		return
	}

	// trans
	userOutList := make([]*model.UserInfo, 0)
	for _, uID := range userIDs {
		userInput := idl.TransUserInput{
			CurUser:  curUser,
			User:     userMap[uID],
			UserStat: nil,
			IsFollow: 0,
			IsFans:   0,
		}
		userInfo := idl.TransUser(&userInput)
		userOutList = append(userOutList, userInfo)
	}

	handler.SendResponse(c, errno.OK, ListResponse{
		TotalCount: 0,
		HasMore:    hasMore,
		PageKey:    "last_id",
		PageValue:  pageValue,
		Items:      userOutList,
	})
}
