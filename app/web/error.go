package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Error404(c *gin.Context) {
	c.HTML(http.StatusOK, "error/404", gin.H{
		"title": "404未找到",
		"ctx":   c,
	})
}
