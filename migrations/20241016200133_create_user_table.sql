-- +goose Up
CREATE TABLE IF NOT EXISTS "user"
(
    id         BIGSERIAL PRIMARY KEY,
    info_id    BIGINT REFERENCES user_info (id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);


-- +goose Down
DROP TABLE IF EXISTS "user" CASCADE;