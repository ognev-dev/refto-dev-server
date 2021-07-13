package model

import (
	"context"
	"time"
)

type EntityData map[string]interface{}

type Entity struct {
	ID        int64      `json:"id"`
	RepoID    int64      `json:"repo_id"`
	Path      string     `json:"path"`
	Title     string     `json:"title"`
	Type      string     `json:"type" pg:",use_zero"`
	Data      EntityData `json:"data"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`

	Repository  *Repository  `json:"repository" pg:",fk:repo_id"`
	Topics      []Topic      `json:"-" pg:"-"`
	Collections []Collection `json:"collections" pg:"-"`

	// composite fields
	SourceURL  string `json:"source_url" pg:"-"`
	EditURL    string `json:"edit_url" pg:"-"`
	CommitsURL string `json:"commits_url" pg:"-"`
}

func (m *Entity) AfterScan(ctx context.Context) error {
	if m.Repository != nil {
		path := GithubAddr + m.Repository.Path
		m.SourceURL = path + "/blob/" + m.Repository.DefaultBranch + "/" + m.Path
		m.EditURL = path + "/edit/" + m.Repository.DefaultBranch + "/" + m.Path
		m.CommitsURL = path + "/commits/" + m.Repository.DefaultBranch + "/" + m.Path
	}

	return nil
}

func (m *Entity) BeforeInsert(ctx context.Context) (context.Context, error) {
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}
	return ctx, nil
}
