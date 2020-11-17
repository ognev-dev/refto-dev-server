package model

import (
	"context"
	"time"
)

type Collection struct {
	ID        int64      `json:"id"`
	Token     string     `json:"token"`
	UserID    int64      `json:"user_id"`
	Name      string     `json:"name"`
	Private   bool       `json:"private" pg:",use_zero"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`

	User *User `json:"-"`
}

func (m *Collection) BeforeInsert(ctx context.Context) (context.Context, error) {
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}
	return ctx, nil
}

func (m *Collection) BeforeUpdate(ctx context.Context) (context.Context, error) {
	now := time.Now()
	m.UpdatedAt = &now

	return ctx, nil
}
