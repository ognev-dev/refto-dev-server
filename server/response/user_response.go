package response

import (
	"time"

	"github.com/refto/server/database/model"
)

type LoginWithGithub struct {
	User      model.User `json:"user"`
	Token     string     `json:"token"`
	ExpiresAt time.Time  `json:"expires_at"`
}
