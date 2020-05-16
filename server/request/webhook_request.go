package request

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
	Repo GitHubRepoPushedRepo `json:"repository" binding:"required"`
}

type GitHubRepoPushedRepo struct {
	CloneURL string `json:"clone_url" binding:"required"`
}
