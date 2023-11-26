-- +goose Up
    CREATE TABLE IF NOT EXISTS products (
        id uuid default gen_random_uuid(),
        category_id uuid not null,
        producer varchar(255) not null,
        model varchar(255) not null,
        additional varchar(255) default null,
        operating_system int,
        description text,
        photos text[],
        price double precision default 0,
        old_price double precision default 0
    );
-- +goose Down