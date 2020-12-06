package util

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/nosurf"
	"server/middleware"
)

func HtmlWithContext(c *gin.Context, status int, template string, args gin.H) {
	merged := make(map[string]interface{})
	for k, v := range args {
		merged[k] = v
	}
	user := middleware.CurrentUser(c)

	if user != nil {
		merged["is_authenticated"] = true
		merged["user_id"] = user.ID
		merged["user_name"] = user.Name
	}

	merged["csrf_token"] = nosurf.Token(c.Request)

	c.HTML(status, template, merged)
}
