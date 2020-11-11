package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Error404 return 404 page
func Error404(c *gin.Context) {
	c.HTML(http.StatusOK, "error/404", gin.H{
		"title": "404未找到",
		"ctx":   c,
	})
}
