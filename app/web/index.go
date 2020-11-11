package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Index home page
func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index", gin.H{
		"title": "首页",
		"ctx":   c,
	})
}
