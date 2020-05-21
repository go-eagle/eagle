package handler

import (
	"net/http"
	"os"

	"github.com/1024casts/snake/pkg/errno"

	"github.com/gin-gonic/gin"
)

// Response api的返回结构体
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// SendResponse 返回json
func SendResponse(c *gin.Context, err error, data interface{}) {
	code, message := errno.DecodeErr(err)

	// always return http.StatusOK
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// GetUserID 返回用户id
func GetUserID(c *gin.Context) uint64 {
	if c == nil {
		return 0
	}

	// uid 必须和 middleware/auth 中的 uid 命名一致
	if v, exists := c.Get("uid"); exists {
		uid, ok := v.(uint64)
		if !ok {
			return 0
		}

		return uid
	}
	return 0
}

// RouteNotFound 未找到相关路由
func RouteNotFound(c *gin.Context) {
	c.String(http.StatusNotFound, "the route not found")
}

// getHostname 获取主机名
func getHostname() string {
	name, err := os.Hostname()
	if err != nil {
		name = "unknown"
	}
	return name
}

// healthCheckResponse 健康检查响应结构体
type healthCheckResponse struct {
	Status   string `json:"status"`
	Hostname string `json:"hostname"`
}

// HealthCheck will return OK if the underlying BoltDB is healthy. At least healthy enough for demoing purposes.
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, healthCheckResponse{Status: "UP", Hostname: getHostname()})
}
