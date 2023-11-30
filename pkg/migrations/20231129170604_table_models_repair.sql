-- +goose Up
    CREATE TABLE IF NOT EXISTS models_repair (
        id uuid default gen_random_uuid(),
        category_repair_id uuid not null,
        name varchar(64) not null
    );
-- +goose Down