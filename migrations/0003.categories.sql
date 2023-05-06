CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT public.uuid_generate_v4(),
    userid UUID REFERENCES users(id) NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    updated_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    deleted_at TIMESTAMP
);