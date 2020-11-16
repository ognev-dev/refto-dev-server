package middleware

import (
	"github.com/gin-gonic/gin"
)

func RequestClient(c *gin.Context) {
	name := c.Request.Header.Get("X-Client")

	c.Set("ClientName", name)
	c.Next()
}
