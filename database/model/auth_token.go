package model

import (
	"context"
	"time"

	"github.com/refto/server/config"
)

type AuthToken struct {
	ID         int64
	UserID     int64
	User       *User
	ClientName string
	ClientIP   string
	UserAgent  string
	Token      string
	CreatedAt  time.Time
	ExpiresAt  time.Time
}

func (m *AuthToken) BeforeInsert(ctx context.Context) (context.Context, error) {
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}
	if m.ExpiresAt.IsZero() {
		m.ExpiresAt = time.Now().Add(config.Get().AuthTokenLifetime)
	}
	return ctx, nil
}
