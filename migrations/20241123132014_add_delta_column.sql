-- +goose Up
-- +goose StatementBegin
ALTER TABLE friend
ADD COLUMN delta VARCHAR;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE friend
DROP COLUMN delta;
-- +goose StatementEnd
