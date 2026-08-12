CREATE TABLE labs (
    id UUID PRIMARY KEY,
    slug TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    difficulty TEXT NOT NULL,
    timeout_minutes INTEGER NOT NULL,
    image TEXT NOT NULL,
    cpu_limit TEXT NOT NULL,
    memory_limit TEXT NOT NULL,
    is_published BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_labs_slug
    ON labs(slug);

CREATE INDEX idx_labs_is_published
    ON labs(is_published);