-- +goose Up
-- +goose StatementBegin
ALTER TABLE items DROP COLUMN IF EXISTS location;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE items ADD COLUMN location VARCHAR(255);
-- +goose StatementEnd