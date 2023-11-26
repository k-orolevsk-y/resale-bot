-- +goose Up
    CREATE TABLE IF NOT EXISTS users (
        id bigint unique,
        tag varchar(255) unique default null,
        is_manager bool default false not null,
        registered_at timestamptz default now() not null
    );
-- +goose Down