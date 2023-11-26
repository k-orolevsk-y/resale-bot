-- +goose Up
    CREATE TABLE IF NOT EXISTS reservation (
        id uuid default gen_random_uuid(),
        user_id uuid not null,
        product_id uuid not null,
        created_at timestamptz default now()
    );
-- +goose Down