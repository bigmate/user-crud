-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id         UUID        NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name TEXT                             DEFAULT '',
    last_name  TEXT                             DEFAULT '',
    nickname   TEXT                             DEFAULT '',
    email      TEXT UNIQUE NOT NULL,
    password   TEXT        NOT NULL,
    country    TEXT                             DEFAULT '',
    created_at TIMESTAMP                        DEFAULT NOW(),
    updated_at TIMESTAMP                        DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
