package request

import githubpullrequest "github.com/refto/server/service/github_pull_request"

// docs: https://developer.github.com/webhooks/event-payloads/#example-delivery

type GitHubWebHookHeaders struct {
	// Name of the event that triggered the delivery
	EventName string `header:"X-GitHub-Event"`
	// A GUID to identify the delivery.
	EventID string `header:"X-GitHub-Delivery"`

	// The HMAC hex digest of the response body.
	// This header will be sent if the webhook is configured with a secret.
	// The HMAC hex digest is generated using the sha1 hash function and the secret as the HMAC key.
	EventSig string `header:"X-Hub-Signature"`
}

type GitHubRepoPushed struct {
	Repo        GitHubRepoPushedRepo `json:"repository" binding:"required"`
	PullRequest GitHubPullRequest    `json:"pull_request"`
}

type GitHubRepoPushedRepo struct {
	CloneURL string `json:"clone_url" binding:"required"`
}

type GitHubPullRequestEvent struct {
	Action      githubpullrequest.Action `json:"action"`
	Number      int                      `json:"number"`
	PullRequest GitHubPullRequest        `json:"pull_request"`
}

type GitHubPullRequest struct {
	Title      string          `json:"title"`
	HTMLURL    string          `json:"html_url"`
	User       GitHubUser      `json:"user"`
	CommitsURL string          `json:"commits_url"`
	Head       PullRequestHead `json:"head"`
}

type GitHubUser struct {
	AvatarURL string `json:"avatar_url"`
	Login     string `json:"login"`
	HTMLURL   string `json:"html_url"`
}

type PullRequestHead struct {
	SHA  string              `json:"sha"`
	Repo PullRequestHeadRepo `json:"repo"`
}

type PullRequestHeadRepo struct {
	CloneURL string `json:"clone_url"`
}
