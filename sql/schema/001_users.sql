-- +goose Up
CREATE TABLE "user" (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    email TEXT NOT NULL UNIQUE
);


-- +goose Down
DROP TABLE "user";