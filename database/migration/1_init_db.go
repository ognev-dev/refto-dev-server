package migration

import (
	"github.com/go-pg/migrations/v7"
)

func init() {
	mg := `
CREATE TABLE topics
(
    id         BIGSERIAL PRIMARY KEY NOT NULL,
    name       TEXT UNIQUE           NOT NULL,
    color      TEXT                  NOT NULL,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE entities
(
    id         BIGSERIAL PRIMARY KEY   NOT NULL,
    token      TEXT UNIQUE             NOT NULL,
    name       TEXT                    NOT NULL,
    type       TEXT                    NOT NULL,
    data       JSONB                   NOT NULL,
    created_at TIMESTAMPTZ             NOT NULL,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE entity_topics
(
    entity_id BIGINT NOT NULL REFERENCES entities (id),
    topic_id  BIGINT NOT NULL REFERENCES topics (id)
);
`

	migrations.MustRegister(func(db migrations.DB) (err error) {
		_, err = db.Exec(mg)
		return
	})
}
