package datawarden

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	githubpullrequest "github.com/refto/server/service/github_pull_request"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v32/github"
	"github.com/refto/server/config"
)

const (
	// I'll link commit fails statuses to PR's comment of the failure description
	// 					  "https://github.com/refto/data/pull/1#issuecomment-1"
	PRCommentLinkFormat = "https://github.com/%s/%s/pull/%d#issuecomment-%d"
)

var (
	instOnce sync.Once
	instance *service
)

type service struct {
	repoOwner string
	repoName  string
	github    *github.Client
}

func Service() *service {
	instOnce.Do(func() {
		repo := config.Get().GitHub.DataRepo

		repoOwner, repoName, err := getRepoOwnerAndNameFromRepoAddr(repo)
		if err != nil {
			panic("data warden: service init: " + err.Error())
		}

		conf := config.Get().GitHub.DataWarden
		tr := http.DefaultTransport
		itr, err := ghinstallation.NewKeyFromFile(tr, conf.AppID, conf.InstallID, conf.PEMPath)
		if err != nil {
			panic("data warden: app auth: " + err.Error())
		}

		client := github.NewClient(&http.Client{
			Transport: itr,
		})

		instance = &service{
			repoOwner: repoOwner,
			repoName:  repoName,
			github:    client,
		}
	})

	return instance
}

func (s *service) Comment(issueNum int, body string) (*github.IssueComment, *github.Response, error) {
	return s.github.Issues.CreateComment(context.Background(), s.repoOwner, s.repoName, issueNum, &github.IssueComment{
		Body: &body,
	})
}

func (s *service) Status(ref string, status githubpullrequest.Status, desc string, url *string) (*github.RepoStatus, *github.Response, error) {
	state := string(status)
	statusCtx := "data-warden"
	return s.github.Repositories.CreateStatus(context.Background(), s.repoOwner, s.repoName, ref, &github.RepoStatus{
		State:       &state,
		TargetURL:   url,
		Description: &desc,
		Context:     &statusCtx,
	})
}

func (s *service) PRCommentLink(prID int, commentID int64) string {
	return fmt.Sprintf(PRCommentLinkFormat, s.repoOwner, s.repoName, prID, commentID)
}

// DataCheckSuccessMessage ...
// Because status needs message and having always same message is boring,
// so I want to make it random
func (s *service) DataCheckSuccessMessage() string {
	rand.Seed(time.Now().Unix())
	statuses := []string{
		"You made it right",
		"Looks like you go through it",
		"You are successful",
		"Thank you for valid data",
		"Safe and sound",
		"You've done it right, keep going",
		"You've done it right, don't stop",
		"Data is checked and it is not wrecked",
		"Why we ever check your commits",
		"Looks safe to merge",
		// Feel free to add your own relevant sentence or fix existing
	}

	return statuses[rand.Intn(len(statuses))]
}

// getRepoOwnerAndNameFromRepoAddr returns repo owner and repo name by given repo addr
// expected addr should be in format:
// https://github.com/refto/data.git
// TODO maybe other formats exits for GitHub Repo to get owner & name
func getRepoOwnerAndNameFromRepoAddr(addr string) (repoOwner, repoName string, err error) {
	addr = strings.TrimPrefix(addr, "https://github.com/")
	addr = strings.TrimPrefix(addr, "http://github.com/")
	addr = strings.TrimSuffix(addr, ".git")
	addrParts := strings.Split(addr, "/")

	if len(addrParts) != 2 {
		err = fmt.Errorf("invalid repo addr. Expecting format: 'https://github.com/{owner}/{repo}.git', got: '%s'", addr)
		return
	}

	return addrParts[0], addrParts[1], nil
}
