package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/1024casts/snake/pkg/app"
	"github.com/1024casts/snake/pkg/errcode"
)

// JWT 认证中间件
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the json web token.
		ctx, err := app.ParseRequest(c)
		if err != nil {
			app.NewResponse().Error(c, errcode.ErrInvalidToken)
			c.Abort()
			return
		}

		// set uid to context
		c.Set("uid", ctx.UserID)

		c.Next()
	}
}
