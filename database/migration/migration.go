package migration

import (
	"github.com/go-pg/migrations/v7"
	"github.com/ognev-dev/bits/database"
)

func Migrate() (int64, int64, error) {
	return migrations.Run(database.Conn(), "up")
}
