package app

import (
	"net/http"

	"github.com/go-eagle/eagle/pkg/utils"

	"github.com/go-eagle/eagle/pkg/errcode"

	"github.com/gin-gonic/gin"
)

// Response define a response struct
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Details []string    `json:"details,omitempty"`
}

// NewResponse return a response
func NewResponse() *Response {
	return &Response{}
}

// Success return a success response
func (r *Response) Success(c *gin.Context, data interface{}) {
	if data == nil {
		data = gin.H{}
	}

	c.JSON(http.StatusOK, Response{
		Code:    errcode.Success.Code(),
		Message: errcode.Success.Msg(),
		Data:    data,
	})
}

// Error return a error response
func (r *Response) Error(c *gin.Context, err error) {
	if err == nil {
		c.JSON(http.StatusOK, Response{
			Code:    errcode.Success.Code(),
			Message: errcode.Success.Msg(),
			Data:    gin.H{},
		})
		return
	}

	if v, ok := err.(*errcode.Error); ok {
		response := Response{
			Code:    v.Code(),
			Message: v.Msg(),
			Data:    gin.H{},
			Details: []string{},
		}
		details := v.Details()
		if len(details) > 0 {
			response.Details = details
		}
		c.JSON(v.StatusCode(), response)
		return
	}
}

// RouteNotFound 未找到相关路由
func RouteNotFound(c *gin.Context) {
	c.String(http.StatusNotFound, "the route not found")
}

// healthCheckResponse 健康检查响应结构体
type healthCheckResponse struct {
	Status   string `json:"status"`
	Hostname string `json:"hostname"`
}

// HealthCheck will return OK if the underlying BoltDB is healthy. At least healthy enough for demoing purposes.
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, healthCheckResponse{Status: "UP", Hostname: utils.GetHostname()})
}
