package user

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/go-eagle/eagle/api"
	"github.com/go-eagle/eagle/internal/ecode"
	"github.com/go-eagle/eagle/internal/service"
	"github.com/go-eagle/eagle/pkg/errcode"
	"github.com/go-eagle/eagle/pkg/log"
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

	_, err := service.Svc.Users().GetUserByID(c.Request.Context(), uint64(userID))
	if err != nil {
		api.SendResponse(c, ecode.ErrUserNotFound, nil)
		return
	}

	lastIDStr := c.DefaultQuery("last_id", "0")
	lastID, _ := strconv.Atoi(lastIDStr)
	limit := 10

	userFollowList, err := service.Svc.Relations().GetFollowingUserList(c.Request.Context(), uint64(userID), uint64(lastID), limit+1)
	if err != nil {
		log.Warnf("get following user list err: %+v", err)
		response.Error(c, errcode.ErrInternalServer)
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

	userOutList, err := service.Svc.Users().BatchGetUsers(c.Request.Context(), curUserID, userIDs)
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
