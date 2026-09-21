CREATE TABLE links (
    id BIGSERIAL PRIMARY KEY,
    original_url TEXT NOT NULL,
    code TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL
);
