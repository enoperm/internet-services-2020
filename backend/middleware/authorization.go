package middleware

import (
	"net/http"
	"server/model"

	"github.com/gin-gonic/gin"
)

const (
	ContextUser = "user-model"
)

func WithUser(getLoggedIn func(c *gin.Context) *model.User) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := getLoggedIn(c)
		if c != nil {
			c.Set(ContextUser, user)
		}
		c.Next()
	}
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := CurrentUser(c)
		if user == nil {
			c.HTML(http.StatusUnauthorized, "unauth", gin.H{
				"title": "Unauthorized",
			})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}
}

func CurrentUser(c *gin.Context) *model.User {
	userIf, found := c.Get(ContextUser)
	if found && userIf.(*model.User) != nil {
		return userIf.(*model.User)
	}
	return nil
}
