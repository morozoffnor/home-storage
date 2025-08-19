-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN homes INTEGER [];

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN IF EXISTS homes;
-- +goose StatementEnd