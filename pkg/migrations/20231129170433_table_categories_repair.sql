-- +goose Up
    CREATE TABLE IF NOT EXISTS categories_repair (
        id uuid default gen_random_uuid(),
        name varchar(64) not null
    );
-- +goose Down