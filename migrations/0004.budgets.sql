CREATE TABLE IF NOT EXISTS budgets (
    id UUID PRIMARY KEY DEFAULT public.uuid_generate_v4(),
    userid UUID REFERENCES users(id) NOT NULL,
    name VARCHAR(255) NOT NULL,
    category UUID REFERENCES categories(id) NOT NULL,
    amount INTEGER NOT NULL,
    period TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    updated_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    deleted_at TIMESTAMP
);