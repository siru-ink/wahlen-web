-- +goose Up
CREATE TABLE roles (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

-- +goose StatementBegin
INSERT INTO roles (name) VALUES ('administrator');
INSERT INTO roles (name) VALUES ('counter');
-- +goose StatementEnd

CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    access INTEGER NOT NULL REFERENCES roles(id)
);

CREATE TABLE voters (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    random_key TEXT NOT NULL UNIQUE,
    has_voted INTEGER NOT NULL DEFAULT FALSE
);

CREATE TABLE votes (
    id INTEGER PRIMARY KEY,
    user_key TEXT NOT NULL UNIQUE REFERENCES voters(random_key),
    voting_result TEXT NOT NULL
);

-- +goose Down
DROP TABLE users;
DROP TABLE roles;
DROP TABLE votes;
DROP TABLE voters;
