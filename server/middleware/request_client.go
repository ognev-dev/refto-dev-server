package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/refto/server/server/request"
)

func RequestClient(c *gin.Context) {
	request.SetClient(c)
	c.Next()
}
