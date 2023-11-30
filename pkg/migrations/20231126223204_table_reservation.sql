-- +goose Up
    CREATE TABLE IF NOT EXISTS reservation (
        id uuid default gen_random_uuid(),
        user_id bigint not null,
        product_id uuid not null unique,
        created_at timestamptz default now(),
        completed int default 0 not null
    );
-- +goose Down