package handler

import (
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/go-git/go-git/v5"
	"github.com/refto/server/config"
	"github.com/refto/server/server/request"
	dataimport "github.com/refto/server/service/data_import"
	datawarden "github.com/refto/server/service/data_warden"
	githubpullrequest "github.com/refto/server/service/github_pull_request"
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
	if headers.EventName != githubwebhook.PushEvent {
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
		log.Error("[ERROR] github's webhook (push) invalid signature")
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
	// TODO: should make selective import using diff
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

		_, err = jsonschema.Validate(conf.Dir.Data)
		if err != nil {
			log.Error("[ERROR] data validate: " + err.Error())
			return
		}

		err = dataimport.Import()
		if err != nil {
			log.Error("[ERROR] data validate: " + err.Error())
			return
		}

		log.Info("Data import from repository completed")
	}()
	c.JSON(http.StatusOK, "ok")
}

func ProcessPullRequestActions(c *gin.Context) {
	var headers request.GitHubWebHookHeaders
	err := c.ShouldBindHeader(&headers)
	if err != nil {
		Abort(c, err)
		return
	}

	// event must be a pull_request
	if headers.EventName != githubwebhook.PullRequestEvent {
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

	conf := config.Get()
	validSig, err := githubwebhook.ValidMAC(body, headers.EventSig, conf.GitHub.DataPushedHookSecret)
	if err != nil {
		log.Error("[ERROR] " + err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if !validSig {
		log.Error("[ERROR] github's webhook (pull_request) invalid signature")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var req request.GitHubPullRequestEvent
	err = c.ShouldBindJSON(&req)
	if err != nil {
		Abort(c, err)
		return
	}

	if !req.Action.ShouldValidate() {
		c.JSON(http.StatusOK, "ok")
		return
	}

	// TODO good to check for merge conflicts here
	// but I don't know how to, seems like github API didnt have info on this

	// because data check might take some time, mark HEAD as status "pending"
	// (in fact it is "in-progress", but some feedback better than nothing)
	_, _, err = datawarden.Service().Status(
		req.PullRequest.Head.SHA,
		githubpullrequest.StatusPending,
		"Checking data...",
		nil,
	)
	if err != nil {
		log.Error("[ERROR] data warden set commit status: " + err.Error())
		return
	}

	go func() {
		var err error

		// if any errors happened, i'd mark the HEAD as failure and comment error to pull request
		// the error might be not related to pull request checks (like internal error or whatever)
		defer func() {
			if err != nil {
				comment, _, err2 := datawarden.Service().Comment(req.Number, err.Error())
				if err2 != nil {
					log.Error("[ERROR] data warden comment: " + err2.Error())
					return
				}

				var commentURL *string
				if comment.ID != nil {
					url := datawarden.Service().PRCommentLink(req.Number, *comment.ID)
					commentURL = &url
				}

				_, _, err = datawarden.Service().Status(
					req.PullRequest.Head.SHA,
					githubpullrequest.StatusFailure,
					err.Error(),
					commentURL,
				)
				if err != nil {
					log.Error("[ERROR] data warden set commit status: " + err.Error())
					return
				}
			}
		}()

		// to make sure that data checks will not go into conflicts
		// i'd clone each PR's HEAD in "pr_{PR_ID}_{HEAD_SHA}" directory
		cloneDir := fmt.Sprintf("pr_%d_%s", req.Number, req.PullRequest.Head.SHA)
		_ = os.MkdirAll(path.Join("pr-checks", cloneDir), 0755)

		_, err = git.PlainClone(cloneDir, false, &git.CloneOptions{
			URL:               req.PullRequest.Head.Repo.CloneURL,
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		})
		if err != nil {
			err = fmt.Errorf("[ERROR] git clone: " + err.Error())
			log.Error(err)
			return
		}

		defer func() {
			err := os.RemoveAll(cloneDir)
			if err != nil {
				err = fmt.Errorf("[ERROR] remove cloned repo (%s) %s: ", cloneDir, err.Error())
				log.Error(err)
			}
		}()

		// TODO validate not only json schema but everything that should be validated
		// like URLs must be valid URL, dates must be valid dates
		// and so on (probably can be done with json schema custom validators)
		_, err = jsonschema.Validate(cloneDir)
		if err != nil {
			err = fmt.Errorf("[ERROR] data validate: %s", err.Error())
			log.Error("[ERROR] data validate: " + err.Error())
			return
		}

		// TODO make test run on database
		// 1. copy current database
		// 2. run import data into it
		// 3. check if any errs
		// 4. delete this (copied) database

		// TODO set reviewers according to topics
		//_, _, err = client.PullRequests.RequestReviewers(context.Background(), "refto", "data", 1, github.ReviewersRequest{
		//	Reviewers: []string{
		//		"data-warden",
		//	},
		//})

		// all good, mark commit as success
		_, _, err = datawarden.Service().Status(
			req.PullRequest.Head.SHA,
			githubpullrequest.StatusSuccess,
			datawarden.Service().DataCheckSuccessMessage(),
			nil,
		)
		if err != nil {
			log.Error("[ERROR] data warden set commit status: " + err.Error())
			return
		}
	}()

	c.JSON(http.StatusOK, "ok")
}
