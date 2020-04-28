package migrations

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
);

CREATE TABLE data
(
    id         BIGSERIAL PRIMARY KEY   NOT NULL,
    token      TEXT UNIQUE PRIMARY KEY NOT NULL,
    name       TEXT                    NOT NULL,
    type       INT                     NOT NULL,
    data       JSONB                   NOT NULL,
    created_at TIMESTAMPTZ             NOT NULL,
    updated_at TIMESTAMPTZ,
);

CREATE TABLE data_topics
(
    data_id   BIGINT NOT NULL REFERENCES data (id),
    topics_id BIGINT NOT NULL REFERENCES topics (id)
);
`

	migrations.MustRegister(func(db migrations.DB) (err error) {
		_, err = db.Exec(mg)
		return
	})
}
