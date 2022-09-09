package middleware

import (
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/errcode"
)

// Auth authorize user
func Auth(paths ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ignore some path
		// eg: register, login, logout
		if len(paths) > 0 {
			path := c.Request.URL.Path
			pathsStr := strings.Join(paths, "|")
			reg := regexp.MustCompile("(" + pathsStr + ")")
			if reg.MatchString(path) {
				return
			}
		}

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
