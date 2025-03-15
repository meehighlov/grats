-- +goose Up
-- +goose StatementBegin
ALTER TABLE chat ADD COLUMN silent_notifications INTEGER DEFAULT 1;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE chat DROP COLUMN silent_notifications;
-- +goose StatementEnd
