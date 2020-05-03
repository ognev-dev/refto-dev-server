package model

import "time"

type Topic struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Color     string `json:"color"`
	DeletedAt *time.Time
}
