CREATE TABLE IF NOT EXISTS recurring_transactions (
    id UUID PRIMARY KEY DEFAULT public.uuid_generate_v4(),
    userid UUID REFERENCES users(id) NOT NULL,
    amount INTEGER NOT NULL,
    category UUID REFERENCES categories(id) NOT NULL,
    description VARCHAR(255),
    vendor VARCHAR(255) NOT NULL,
    date TIMESTAMP NOT NULL,
    type VARCHAR(255) NOT NULL,
    last_execution TIMESTAMP NOT NULL DEFAULT '1970-01-01 00:00:00',
    next_execution TIMESTAMP NOT NULL,
    unit_of_measure VARCHAR(50) NOT NULL,
    frequency_count INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    updated_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    deleted_at TIMESTAMP
);