-- +goose Up
CREATE TABLE IF NOT EXISTS roles
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
);

INSERT INTO roles (name)
VALUES ('UNKNOWN'),('ADMIN'),('USER');

-- +goose Down
DROP TABLE IF EXISTS roles;