package request

import (
	"github.com/gin-gonic/gin"
	"github.com/refto/server/database/model"
)

type Pagination struct {
	Page  int `json:"page,omitempty" form:"page"`
	Limit int `json:"limit,omitempty" form:"per_page"`
}

type NoValidation struct{}

func (r *NoValidation) Validate(*gin.Context) (err error) {
	return
}

const ctxUserKey = "request_user"
const ctxClientKey = "request_client"

func User(c *gin.Context) model.User {
	u := c.MustGet(ctxUserKey)
	return u.(model.User)
}

func SetUser(c *gin.Context, u model.User) {
	c.Set(ctxUserKey, u)
}

func Client(c *gin.Context) string {
	return c.GetString(ctxClientKey)
}

func SetClient(c *gin.Context) {
	name := c.Request.Header.Get("X-Client")
	c.Set(ctxClientKey, name)
}
