-- +goose Up
CREATE TABLE IF NOT EXISTS role
(
    id INTEGER PRIMARY KEY CHECK (id >= 0),
    name VARCHAR(50) NOT NULL UNIQUE
);

INSERT INTO role (id, name) VALUES (0, 'UNKNOWN'), (1, 'ADMIN'), (2, 'USER');

-- +goose Down
DROP TABLE IF EXISTS role;