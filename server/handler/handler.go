package handler

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"
	"github.com/ognev-dev/bits/config"
	serverError "github.com/ognev-dev/bits/server/error"
	"github.com/ognev-dev/bits/server/response"
	log "github.com/sirupsen/logrus"
)

type Validatable interface {
	Validate() error
}

func abort400(c *gin.Context, err error) {
	Abort(c, serverError.New400(err.Error()))
}

func abort422(c *gin.Context, err error) {
	Abort(c, serverError.New422(err.Error()))
}

func Abort(c *gin.Context, err error) {
	resp := response.Error{}
	code := http.StatusInternalServerError

	switch err.(type) {
	case serverError.Error:
		e := err.(serverError.Error)
		code = e.Code
		resp.Error = e.Error()
	case serverError.List:
		resp.Errors = err.(serverError.List)
	case serverError.Input:
		resp.InputErrors = err.(serverError.Input)
	default:
		resp.Error = err.Error()
		switch err {
		case pg.ErrNoRows:
			code = http.StatusNotFound
		}
	}

	if !config.IsReleaseEnv() {
		log.Println(err.Error())
		log.Println(string(debug.Stack()))
	}

	c.AbortWithStatusJSON(code, resp)
}

func bindJSON(c *gin.Context, req Validatable) (ok bool) {
	err := c.ShouldBindJSON(req)
	if err != nil {
		abort400(c, err)
		return
	}

	err = req.Validate()
	if err != nil {
		abort422(c, err)
		return
	}

	return true
}

func bindQuery(c *gin.Context, req Validatable) (ok bool) {
	err := c.ShouldBindQuery(req)
	if err != nil {
		abort400(c, err)
		return
	}

	err = req.Validate()
	if err != nil {
		abort422(c, err)
		return
	}

	return true
}
