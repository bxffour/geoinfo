CREATE TABLE IF NOT EXISTS countries (
    id SERIAL PRIMARY KEY,
    country JSONB,
    version integer NOT NULL DEFAULT 1
);