-- +goose Up
    CREATE TABLE IF NOT EXISTS repair (
        id uuid default gen_random_uuid(),
        producer_name varchar(64) not null,
        model_name varchar(64) not null,
        name varchar(64) not null,
        description text,
        price double precision not null default 0.0
    );
-- +goose Down