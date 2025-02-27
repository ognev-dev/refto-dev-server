package test

import (
	"net/http"
	"testing"

	"github.com/refto/server/server/request"
	. "github.com/refto/server/test/apitest"
)

func TestWebhookDataPushed(t *testing.T) {
	// TODO
	// mock repo so this test will work)
	t.Skip()

	//conf := config.Get()
	req := request.GitHubRepoPushed{
		Repo: request.GitHubRepo{
			//CloneURL: conf.GitHub.DataRepo,
		},
	}

	//sig, err := githubwebhook.MakeHMAC(
	//	`{"repository":{"clone_url":"`+conf.GitHub.DataRepo+`"}}`,
	//	conf.GitHub.DataPushedHookSecret,
	//)
	//assert.NotError(t, err)

	POST(t, Request{
		Path:         "hooks/data-pushed",
		Body:         req,
		BindResponse: nil,
		AssertStatus: http.StatusOK,
		Headers: Headers{
			"X-GitHub-Event": "push",
			//"X-Hub-Signature": sig,
		},
	})
}

func TestWebhookPullRequest(t *testing.T) {
	t.Skip()

	//conf := config.Get()
	req := request.GitHubRepoPushed{
		Repo: request.GitHubRepo{
			//CloneURL: conf.GitHub.DataRepo,
		},
	}

	//sig, err := githubwebhook.MakeHMAC(
	//`{"repository":{"clone_url":"`+conf.GitHub.DataRepo+`"}}`,
	//conf.GitHub.DataPushedHookSecret,
	//)
	//assert.NotError(t, err)

	POST(t, Request{
		Path:         "hooks/data-pushed",
		Body:         req,
		BindResponse: nil,
		AssertStatus: http.StatusOK,
		Headers: Headers{
			"X-GitHub-Event": "push",
			//"X-Hub-Signature": sig,
		},
	})
}
