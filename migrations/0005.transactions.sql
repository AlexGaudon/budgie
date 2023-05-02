CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT public.uuid_generate_v4(),
    userid UUID REFERENCES users(id) NOT NULL,
    amount INTEGER NOT NULL,
    category UUID REFERENCES categories(id) NOT NULL,
    description VARCHAR(255),
    vendor VARCHAR(255) NOT NULL,
    date TIMESTAMP NOT NULL,
    type VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    updated_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    deleted_at TIMESTAMP
)