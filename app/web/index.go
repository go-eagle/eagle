package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {

	c.HTML(http.StatusOK, "index", gin.H{
		"title": "首页",
		"ctx":   c,
	})
}
