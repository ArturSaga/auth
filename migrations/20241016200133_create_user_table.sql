-- +goose Up
CREATE TYPE user_role AS ENUM ('UNKNOWN', 'ADMIN', 'USER');
CREATE TABLE IF NOT EXISTS users
(
    id               BIGSERIAL PRIMARY KEY,
    name             VARCHAR(255) NOT NULL,
    email            VARCHAR(255) NOT NULL UNIQUE,
    password_hash    VARCHAR(60) NOT NULL,
    role             user_role NOT NULL DEFAULT 'UNKNOWN',
    created_at       TIMESTAMPTZ DEFAULT NOW(),
    updated_at       TIMESTAMPTZ DEFAULT NOW()
);

-- +goose Down
DROP TYPE IF EXISTS user_role;
DROP TABLE IF EXISTS users CASCADE;
