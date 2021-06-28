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
	})
}

func (r *Response) Error(c *gin.Context, err *errcode.Error) {
	response := gin.H{"code": err.Code(), "msg": err.Msg(), "data": gin.H{}}
	details := err.Details()
	if len(details) > 0 {
		response["details"] = details
	}

	c.JSON(err.StatusCode(), response)
}
