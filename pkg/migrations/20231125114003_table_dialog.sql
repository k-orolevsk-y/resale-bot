-- +goose Up
    CREATE TABLE IF NOT EXISTS dialog (
        id uuid default gen_random_uuid(),
        user_id bigint not null,
        manager_id bigint not null,
        started_at timestamptz default now() not null,
        ended_at timestamptz default null
    );
-- +goose Down