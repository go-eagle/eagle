package app

import (
	"net/http"

	"github.com/1024casts/snake/pkg/errcode"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Details []string    `json:"details"`
}

func NewResponse() *Response {
	return &Response{}
}

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
	if err != nil {
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

	c.JSON(http.StatusOK, Response{
		Code:    errcode.Success.Code(),
		Message: errcode.Success.Msg(),
		Data:    gin.H{},
	})
}
