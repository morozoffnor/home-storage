-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS homes (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

ALTER TABLE containers ADD COLUMN home_id INTEGER REFERENCES homes(id) ON DELETE SET NULL;
CREATE INDEX idx_containers_home_id ON containers(home_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_containers_home_id;
ALTER TABLE containers DROP COLUMN IF EXISTS home_id;
DROP TABLE IF EXISTS homes;
-- +goose StatementEnd