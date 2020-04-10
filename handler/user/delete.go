package user

import (
	"github.com/1024casts/snake/handler"

	"github.com/gin-gonic/gin"
)

// Delete 删除用户
// @Summary Delete an user by the user identifier
// @Description Delete user by ID
// @Tags 用户
// @Accept  json
// @Produce  json
// @Param id path uint64 true "The user's database id index num"
// @Success 200 {object} handler.Response "{"code":0,"message":"OK","data":null}"
// @Router /user/{id} [delete]
func Delete(c *gin.Context) {
	//userId, _ := strconv.Atoi(c.Param("id"))
	//if err := model.DeleteUser(uint64(userId)); err != nil {
	//	SendResponse(c, errno.ErrDatabase, nil)
	//	return
	//}

	handler.SendResponse(c, nil, nil)
}
