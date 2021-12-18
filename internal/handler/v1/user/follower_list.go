package user

import (
	"context"
	"strconv"

	"github.com/go-eagle/eagle/internal/service"

	"github.com/go-eagle/eagle/internal/ecode"

	"github.com/gin-gonic/gin"

	"github.com/go-eagle/eagle/pkg/errcode"
	"github.com/go-eagle/eagle/pkg/log"
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

	curUserID := service.GetUserID(c)

	_, err := service.Svc.Users().GetUserByID(c, uint64(userID))
	if err != nil {
		response.Error(c, ecode.ErrUserNotFound)
		return
	}

	lastIDStr := c.DefaultQuery("last_id", "0")
	lastID, _ := strconv.Atoi(lastIDStr)
	limit := 10

	userFollowerList, err := service.Svc.Relations().GetFollowerUserList(context.TODO(), uint64(userID), uint64(lastID), limit+1)
	if err != nil {
		log.Warnf("get follower user list err: %+v", err)
		response.Error(c, errcode.ErrInternalServer)
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

	userOutList, err := service.Svc.Users().BatchGetUsers(context.TODO(), curUserID, userIDs)
	if err != nil {
		log.Warnf("batch get users err: %v", err)
		response.Error(c, errcode.ErrInternalServer)
		return
	}

	response.Success(c, ListResponse{
		TotalCount: 0,
		HasMore:    hasMore,
		PageKey:    "last_id",
		PageValue:  pageValue,
		Items:      userOutList,
	})
}
