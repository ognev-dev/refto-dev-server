package response

import (
	"fmt"

	"github.com/refto/server/database/model"
)

type CreateRepository struct {
	WebhookCreateURL  string `json:"webhook_create_url"`
	WebhookPayloadURL string `json:"webhook_payload_url"`
	WebhookSecret     string `json:"webhook_secret"`
	Path              string `json:"repo_path"`
}

func NewCreateRepository(rep model.Repository) CreateRepository {
	return CreateRepository{
		WebhookCreateURL:  fmt.Sprintf("https://github.com/%s/settings/hooks/new", rep.Path),
		WebhookPayloadURL: "https://refto.dev/api/hooks/data-pushed/",
		WebhookSecret:     rep.Secret,
		Path:              rep.Path,
	}
}

type FilterRepositories struct {
	Data  []model.Repository `json:"data"`
	Count int                `json:"count"`
}
