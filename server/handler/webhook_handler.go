package handler

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-git/go-git/v5"
	"github.com/refto/server/config"
	"github.com/refto/server/server/request"
	dataimport "github.com/refto/server/service/data_import"
	githubwebhook "github.com/refto/server/service/github_webhook"
	jsonschema "github.com/refto/server/service/json_schema"
	log "github.com/sirupsen/logrus"
)

func ImportDataFromRepoByGitHubWebHook(c *gin.Context) {
	conf := config.Get()
	var headers request.GitHubWebHookHeaders
	err := c.ShouldBindHeader(&headers)
	if err != nil {
		Abort(c, err)
		return
	}

	// event must be a push event
	if headers.EventName != "push" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// check signature
	body, err := copyRequestBody(c)
	if err != nil {
		log.Error(err)
		Abort(c, err)
		return
	}

	validSig, err := githubwebhook.ValidMAC(body, headers.EventSig, conf.GitHub.DataPushedHookSecret)
	if err != nil {
		log.Error("[ERROR] " + err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if !validSig {
		log.Error("github's data pushed webhook invalid signature")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var req request.GitHubRepoPushed
	err = c.ShouldBindJSON(&req)
	if err != nil {
		Abort(c, err)
		return
	}

	if req.Repo.CloneURL != conf.GitHub.DataRepo {
		log.Errorf("clone repo (%s) is not same as data repo (%s)", req.Repo.CloneURL, config.Get().GitHub.DataRepo)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// import data on goroutine, because it is nothing to do with request
	// TODO: must make selective import using diff
	go func() {
		log.Info("Starting data import from " + conf.GitHub.DataRepo + " to " + conf.Dir.Data)
		err := os.RemoveAll(conf.Dir.Data)
		if err != nil {
			log.Error("[ERROR] os.RemoveAll: " + err.Error())
			return
		}
		_, err = git.PlainClone(conf.Dir.Data, false, &git.CloneOptions{
			URL:               conf.GitHub.DataRepo,
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		})
		if err != nil {
			log.Error("[ERROR] git clone: " + err.Error())
			return
		}

		err = jsonschema.Validate()
		if err != nil {
			log.Error("[ERROR] data validate: " + err.Error())
			return
		}

		err = dataimport.Process()
		if err != nil {
			log.Error("[ERROR] data validate: " + err.Error())
			return
		}

		log.Info("Data import from repository completed")
	}()
	c.JSON(http.StatusOK, "ok")
}
