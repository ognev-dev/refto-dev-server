package request

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type LoginWithGithub struct {
	Code string `json:"code"`
}

func (r *LoginWithGithub) Validate(*gin.Context) (err error) {
	if r.Code == "" {
		err = errors.New("login code missing")
		return
	}

	return
}
