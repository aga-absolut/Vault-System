-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (         
    username      TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS records (
    record_id   SERIAL PRIMARY KEY,
    record_type TEXT NOT NULL,
    record_meta TEXT UNIQUE NOT NULL,
    record_data BYTEA NOT NULL,
    username    TEXT NOT NULL REFERENCES users(username) ON DELETE CASCADE
); 
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS records CASCADE;
DROP TABLE IF EXISTS users CASCADE;
-- +goose StatementEnd