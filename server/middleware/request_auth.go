package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	se "github.com/refto/server/server/error"
	"github.com/refto/server/server/handler"
)

// RequestAuth middleware
// See RequestUser middleware for auth details
func RequestAuth(c *gin.Context) {
	if !c.GetBool(authSuccessKey) {
		errText := c.GetString(authErrorKey)
		if errText == "" {
			errText = http.StatusText(http.StatusUnauthorized)
		}

		handler.Abort(c, se.New401(errText))
		return
	}

	c.Next()
}
