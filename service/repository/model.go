package repository

import "github.com/refto/server/database/model"

type FilterParams struct {
	Page      int
	Limit     int
	Path      string
	Name      string
	Types     []model.RepoType
	UserID    int64
	Confirmed *bool
}
