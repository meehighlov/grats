-- +goose Up
-- +goose StatementBegin
ALTER TABLE chat ADD COLUMN greeting_template VARCHAR DEFAULT '🔔Сегодня день рождения у %s🥳';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE chat DROP COLUMN greeting_template;
-- +goose StatementEnd 