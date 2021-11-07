package app

import (
	"net/http"

	"github.com/go-eagle/eagle/pkg/errcode"

	"github.com/gin-gonic/gin"
)

// Response define a response struct
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Details []string    `json:"details"`
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
		Details: []string{},
	})
}

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
