package migration

import (
	"github.com/go-pg/migrations/v7"
)

func init() {
	mg := `
CREATE TABLE collections
(
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT      NOT NULL REFERENCES users (id),
    token      TEXT UNIQUE,
    name       TEXT        NOT NULL,
    private    BOOLEAN,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ
);

CREATE TABLE collection_entities
(
    collection_id BIGINT NOT NULL REFERENCES collections (id),
    entity_id     BIGINT NOT NULL REFERENCES entities (id)
);
`

	migrations.MustRegister(func(db migrations.DB) (err error) {
		_, err = db.Exec(mg)
		return
	})
}
