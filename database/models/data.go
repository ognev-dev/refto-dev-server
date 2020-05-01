package models

import (
	"context"
	"time"
)

type Data struct {
	tableName struct{}   `pg:"data"`
	ID        string     `json:"-"`
	Token     string     `json:"token"`
	Name      string     `json:"-"`
	Type      string     `json:"type"`
	Data      string     `json:"data"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func (m *Data) BeforeInsert(ctx context.Context) (context.Context, error) {
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}
	return ctx, nil
}

func (m *Data) BeforeUpdate(ctx context.Context) (context.Context, error) {
	if m.UpdatedAt.IsZero() {
		now := time.Now()
		m.UpdatedAt = &now
	}
	return ctx, nil
}
