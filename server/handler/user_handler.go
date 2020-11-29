package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/refto/server/database/model"
	"github.com/refto/server/server/request"
	"github.com/refto/server/server/response"
	authtoken "github.com/refto/server/service/auth_token"
	"github.com/refto/server/service/github"
	"github.com/refto/server/service/user"
)

func LoginWithGithub(c *gin.Context) {
	var req request.LoginWithGithub
	if !bindJSON(c, &req) {
		return
	}

	ghToken, err := github.GetAccessToken(req.Code)
	if err != nil {
		Abort(c, err)
		return
	}

	gu, err := github.GetUser(ghToken)
	if err != nil {
		Abort(c, err)
		return
	}

	usr, err := user.ResolveFromGithubUser(gu, ghToken)
	if err != nil {
		Abort(c, err)
		return
	}

	token := &model.AuthToken{
		UserID:     usr.ID,
		ClientName: request.Client(c),
		ClientIP:   c.Request.RemoteAddr,
		UserAgent:  c.Request.UserAgent(),
	}
	err = authtoken.Create(token)
	if err != nil {
		Abort(c, err)
		return
	}

	resp := response.LoginWithGithub{
		User:      usr,
		Token:     authtoken.Sign(token),
		ExpiresAt: token.ExpiresAt,
	}

	c.JSON(http.StatusOK, resp)
}
