package handler

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"
	"github.com/refto/server/config"
	serverError "github.com/refto/server/server/error"
	"github.com/refto/server/server/response"
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

func abort401(c *gin.Context, err error) {
	Abort(c, serverError.New422(err.Error()))
}

func Abort(c *gin.Context, err error) {
	resp := response.Error{}
	code := http.StatusInternalServerError

	switch e := err.(type) {
	case serverError.Error:
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

	log.Println(err.Error())
	if !config.IsReleaseEnv() {
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

func copyRequestBody(c *gin.Context) (body []byte, err error) {
	body, err = ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body)) // Write body back

	return
}
