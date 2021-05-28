// GitHub's repository
// https://docs.github.com/en/rest/reference/repos#get-a-repository
// GET https://api.github.com/repos/{user}/{name}

package model

import (
	"context"
	"time"
)

type Repository struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"-"`

	// Path is just a "{user}/{repo}" on GitHub
	Path string `json:"path"`

	Name        string `json:"name"`
	Description string `json:"description"`

	// Secret is needed to authenticate repository from Github
	// It is random string given to the user who creates repository
	// and must be set on repository's push webhook on Github
	// https://github.com/{account}/{repo}/settings/hooks
	// Secret is hashed using bcrypt and available only once
	Secret string `json:"-"`

	Type RepositoryType `json:"type"`

	// Confirmed is a flag to mark that user is confirmed access to repo
	// Confirmed is set to true on first successful import
	Confirmed *bool `json:"-"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func (m *Repository) BeforeInsert(ctx context.Context) (context.Context, error) {
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}

	return ctx, nil
}

func (m *Repository) BeforeUpdate(ctx context.Context) (context.Context, error) {
	now := time.Now()
	m.UpdatedAt = &now

	return ctx, nil
}

type RepositoryType string

const (
	// RepositoryTypePrivate
	// Private repos is available only to user who added it
	RepositoryTypePrivate RepositoryType = "private"

	// RepositoryTypeGlobal
	// Data from global repos is available by default at global level
	RepositoryTypeGlobal RepositoryType = "global"

	// RepositoryTypePublic
	// Public repos will be listed in repos page and in search filters
	RepositoryTypePublic RepositoryType = "public"

	// RepositoryTypeHidden
	// Hidden repos will NOT be listed in repos page and in search filters
	// but can be viewed using link
	RepositoryTypeHidden RepositoryType = "hidden"
)

var RepositoryTypesList = []string{
	"private", "global", "public", "hidden",
}

func (t RepositoryType) IsValid() bool {
	switch t {
	case
		RepositoryTypePrivate,
		RepositoryTypeGlobal,
		RepositoryTypeHidden,
		RepositoryTypePublic:
		return true
	}

	return false
}

func (t RepositoryType) String() string {
	return string(t)
}
