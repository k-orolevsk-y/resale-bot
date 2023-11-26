-- +goose Up
    CREATE TABLE IF NOT EXISTS categories (
        id uuid default gen_random_uuid(),
        name varchar(64) not null,
        c_type int not null
    );
-- +goose Down