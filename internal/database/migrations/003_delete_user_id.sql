-- +goose Up
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_items_user_id;
ALTER TABLE items DROP COLUMN IF EXISTS user_id;
ALTER TABLE items ADD COLUMN container_id INTEGER REFERENCES containers(id) ON DELETE SET NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE items DROP COLUMN IF EXISTS container_id;
ALTER TABLE items ADD COLUMN user_id INTEGER REFERENCES users(id) ON DELETE CASCADE;
CREATE INDEX idx_items_user_id ON items(user_id);
-- +goose StatementEnd