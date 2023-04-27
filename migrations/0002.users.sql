CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT public.uuid_generate_v4(),
    username varchar(255) NOT NULL,
    passwordhash varchar(1024) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    updated_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    deleted_at TIMESTAMP
)