package model

import (
	"context"
	"time"
)

type User struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	Login       string     `json:"login"`
	AvatarURL   string     `json:"avatar_url"`
	TelegramID  int64      `json:"-"`
	GithubID    int64      `json:"-"`
	GithubToken string     `json:"-"`
	Email       string     `json:"-"`
	CreatedAt   time.Time  `json:"created_at"`
	ActiveAt    time.Time  `json:"active_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

func (m *User) BeforeInsert(ctx context.Context) (context.Context, error) {
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}
	return ctx, nil
}

func (m *User) BeforeUpdate(ctx context.Context) (context.Context, error) {
	now := time.Now()
	m.UpdatedAt = &now

	return ctx, nil
}
