-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS "friend";
ALTER TABLE "user" DROP COLUMN IF EXISTS birthday;
ALTER TABLE "chat" DROP COLUMN IF EXISTS greeting_template;
ALTER TABLE "chat" DROP COLUMN IF EXISTS silent_notifications;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "friend" (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255),
    birthday VARCHAR(255),
    chat_id VARCHAR(255),
    user_id VARCHAR(36),
    notify_at VARCHAR(255),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

ALTER TABLE "user" ADD COLUMN birthday VARCHAR(255);
ALTER TABLE "chat" ADD COLUMN greeting_template VARCHAR(255) DEFAULT '🔔Сегодня день рождения у %s🥳';
ALTER TABLE "chat" ADD COLUMN silent_notifications BOOLEAN DEFAULT TRUE;
-- +goose StatementEnd 