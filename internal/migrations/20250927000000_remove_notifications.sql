-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS friend;
DROP TABLE IF EXISTS chat;
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
CREATE TABLE IF NOT EXISTS "chat" (
    id VARCHAR(36) PRIMARY KEY,
    chat_id VARCHAR(255),
    chat_type VARCHAR(50),
    bot_invited_by_id VARCHAR(255),
    created_at TIMESTAMP,
    updated_at TIMESTAMP, 
    greeting_template VARCHAR(255) DEFAULT '🔔Сегодня день рождения у %s🥳', 
    silent_notifications BOOLEAN DEFAULT TRUE,
    UNIQUE(chat_id)
);
-- +goose StatementEnd
