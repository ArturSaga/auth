-- +goose Up
CREATE TABLE IF NOT EXISTS user_info
(
    id               BIGSERIAL PRIMARY KEY,
    name             VARCHAR(255) NOT NULL,
    email            VARCHAR(255) NOT NULL UNIQUE,
    password         VARCHAR(255) NOT NULL,
    password_confirm VARCHAR(255) NOT NULL,
    role_id          INT REFERENCES role (id)
);


-- +goose Down
DROP TABLE IF EXISTS user_info CASCADE;