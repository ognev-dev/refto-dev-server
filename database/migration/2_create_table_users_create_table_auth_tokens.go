package migration

import (
	"github.com/go-pg/migrations/v7"
)

//{
//"access_token": "a15d631eb2e048d2935bc4d6bf2226d6307f3d3a",
//"token_type": "bearer",
//"scope": ""
//}

func init() {
	mg := `
CREATE TABLE users
(
    id           BIGSERIAL PRIMARY KEY,
    telegram_id  INT,
    github_id    BIGINT UNIQUE NOT NULL,
    github_token TEXT,
    avatar_url   TEXT,
    name         TEXT,
    email        TEXT,
    created_at   TIMESTAMPTZ NOT NULL,
    active_at    TIMESTAMPTZ NOT NULL,
    updated_at   TIMESTAMPTZ
);

CREATE TABLE auth_tokens
(
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT      NOT NULL REFERENCES users (id),
    client_ip   TEXT,
    client_name TEXT,
    user_agent  TEXT,
    token       TEXT,
    created_at  TIMESTAMPTZ NOT NULL,
    expires_at  TIMESTAMPTZ NOT NULL
);
`

	migrations.MustRegister(func(db migrations.DB) (err error) {
		_, err = db.Exec(mg)
		return
	})
}
