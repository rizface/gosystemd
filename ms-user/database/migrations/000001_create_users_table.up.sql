create table if not exists users (
    id uuid primary key,
    name varchar not null,
    username varchar not null unique,
    password varchar not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
)