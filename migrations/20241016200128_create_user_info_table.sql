-- +goose Up
CREATE TABLE user_info IF NOT EXISTS
(
    id               BIGSERIAL PRIMARY KEY,
    name             VARCHAR(255) NOT NULL,
    email            VARCHAR(255) NOT NULL UNIQUE,
    password         VARCHAR(255) NOT NULL,
    password_confirm VARCHAR(255) NOT NULL,
    role_id          INT REFERENCES roles (id)
);


-- +goose Down
DROP TABLE IF EXISTS user_info CASCADE;