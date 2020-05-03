package model

import (
	"context"
	"time"
)

type Entity struct {
	ID        int64      `json:"id"`
	Token     string     `json:"token"`
	Name      string     `json:"name"`
	Type      string     `json:"type"`
	Data      string     `json:"data"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`

	Topics []Topic `json:"-" pg:"-"`
}

func (m *Entity) BeforeInsert(ctx context.Context) (context.Context, error) {
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}
	return ctx, nil
}

func (m *Entity) BeforeUpdate(ctx context.Context) (context.Context, error) {
	if m.UpdatedAt.IsZero() {
		now := time.Now()
		m.UpdatedAt = &now
	}
	return ctx, nil
}
