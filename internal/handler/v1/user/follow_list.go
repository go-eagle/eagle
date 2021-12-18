package user

import (
	"strconv"

	"github.com/gin-gonic/gin"

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
	// get the underlying request context
	ctx := c.Request.Context()

	// create a done channel to tell the request it's done
	doneChan := make(chan ListResponse)
	// create a err channel
	errChan := make(chan error)

	// here you put the actual work needed for the request
	// and then send the doneChan with the status and body
	// to finish the request by writing the response
	go func() {
		userIDStr := c.Param("id")
		userID, _ := strconv.Atoi(userIDStr)

		curUserID := service.GetUserID(c)
		log.Infof("cur uid: %d", curUserID)

		_, err := service.Svc.Users().GetUserByID(ctx, uint64(userID))
		if err != nil {
			errChan <- ecode.ErrUserNotFound
			return
		}

		lastIDStr := c.DefaultQuery("last_id", "0")
		lastID, _ := strconv.Atoi(lastIDStr)
		limit := 10

		userFollowList, err := service.Svc.Relations().GetFollowingUserList(ctx, uint64(userID), uint64(lastID), limit+1)
		if err != nil {
			log.Warnf("get following user list err: %+v", err)
			errChan <- errcode.ErrInternalServer
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

		userOutList, err := service.Svc.Users().BatchGetUsers(ctx, curUserID, userIDs)
		if err != nil {
			log.Warnf("batch get users err: %v", err)
			errChan <- errcode.ErrInternalServer
			return
		}

		doneChan <- ListResponse{
			TotalCount: 0,
			HasMore:    hasMore,
			PageKey:    "last_id",
			PageValue:  pageValue,
			Items:      userOutList,
		}
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
