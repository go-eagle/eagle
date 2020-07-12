package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"

	"github.com/1024casts/snake/handler"
)

// old style see: http://researchlab.github.io/2016/03/29/gin-setcookie/
func CookieMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := handler.GetCookieSession(c)
		log.Infof("[middleware] current session: %v", session)
		userId, ok := session.Values["user_id"]
		log.Infof("[middleware] current user_id: %d from cookie", userId)
		if !ok {
			log.Warnf("[middleware] current user_id is not ok")
			c.Redirect(http.StatusMovedPermanently, "/login")
			c.Abort()
			return
		}

		handler.SetLoginCookie(c, userId.(uint64))

		c.Next()
	}
}
