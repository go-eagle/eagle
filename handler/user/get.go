package user

import (
	. "go-api/handler"
	"go-api/model"
	"go-api/pkg/errno"

	"github.com/gin-gonic/gin"
)

// @Summary Get an user by the user identifier
// @Description Get an user by username
// @Tags user
// @Accept  json
// @Produce  json
// @Param username path string true "Username"
// @Success 200 {object} model.UserModel "{"code":0,"message":"OK","data":{"username":"kong","password":"$2a$10$E0kwtmtLZbwW/bDQ8qI8e.eHPqhQOW9tvjwpyo/p05f/f4Qvr3OmS"}}"
// @Router /user/{username} [get]
func Get(c *gin.Context) {
	username := c.Param("username")
	// Get the user by the `username` from the database.
	user, err := model.GetUser(username)
	if err != nil {
		SendResponse(c, errno.ErrUserNotFound, nil)
		return
	}

	SendResponse(c, nil, user)
}
