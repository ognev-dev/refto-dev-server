package middleware

import (
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/refto/server/server/handler"
	authtoken "github.com/refto/server/service/auth_token"
)

func AuthToken(c *gin.Context) {
	sig := c.Request.Header.Get("Authorization")
	if sig == "" {
		handler.Abort(c, errors.New("auth token missing"))
		return
	}

	sig = strings.TrimPrefix(sig, "Bearer ")
	if len(sig) < 128 {
		handler.Abort(c, errors.New("invalid auth token"))
		return
	}

	token := sig[:64]
	elem, err := authtoken.FindByToken(token)
	if err != nil {
		handler.Abort(c, err)
		return
	}

	elem.ClientName = c.GetString("ClientName")
	elem.ClientIP = c.Request.RemoteAddr
	elem.UserAgent = c.Request.UserAgent()

	if sig != authtoken.Sign(&elem) {
		handler.Abort(c, errors.New("invalid auth token signature"))
		return
	}

	if elem.ExpiresAt.Before(time.Now()) {
		handler.Abort(c, errors.New("auth token expired"))
		return
	}

	c.Set("user", *elem.User)
	c.Next()
}
