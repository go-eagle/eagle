package app

import (
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/go-eagle/eagle/pkg/utils"
	"github.com/spf13/cast"
	"google.golang.org/grpc/status"

	"github.com/go-eagle/eagle/pkg/errcode"

	"github.com/gin-gonic/gin"
)

var resp *Response

func init() {
	resp = NewResponse()
}

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
func Success(c *gin.Context, data interface{}) { resp.Success(c, data) }
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
func Error(c *gin.Context, err error) { resp.Error(c, err) }
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
		c.JSON(errcode.ToHTTPStatusCode(v.Code()), response)
		return
	} else {
		// receive gRPC error
		if st, ok := status.FromError(err); ok {
			response := Response{
				Code:    int(st.Code()),
				Message: st.Message(),
				Data:    gin.H{},
				Details: []string{},
			}
			details := st.Details()
			if len(details) > 0 {
				for _, v := range details {
					response.Details = append(response.Details, cast.ToString(v))
				}
			}
			// https://httpstatus.in/
			// https://github.com/grpc-ecosystem/grpc-gateway/blob/master/runtime/errors.go#L15
			// https://github.com/googleapis/googleapis/blob/master/google/rpc/code.proto
			c.JSON(runtime.HTTPStatusFromCode(st.Code()), response)
			return
		}
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

// HealthCheck will return OK if the underlying BoltDB is healthy.
// At least healthy enough for demoing purposes.
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, healthCheckResponse{Status: "UP", Hostname: utils.GetHostname()})
}
