-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "user" (
    id VARCHAR(36) PRIMARY KEY,
    tg_id VARCHAR(255),
    name VARCHAR(255),
    tg_username VARCHAR(255),
    chat_id VARCHAR(255),
    is_admin BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    UNIQUE(tg_id)
);

CREATE TABLE IF NOT EXISTS "chat" (
    id VARCHAR(36) PRIMARY KEY,
    chat_id VARCHAR(255),
    chat_type VARCHAR(50),
    bot_invited_by_id VARCHAR(255),
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    UNIQUE(chat_id)
);

CREATE TABLE IF NOT EXISTS wish (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255),
    chat_id VARCHAR(255),
    user_id VARCHAR(36),
    link VARCHAR(255),
    executor_id VARCHAR(36),
    price VARCHAR(255),
    wish_list_id VARCHAR(36),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS wish_list (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255),
    user_id VARCHAR(36),
    chat_id VARCHAR(255),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS wish_list;
DROP TABLE IF EXISTS wish;
DROP TABLE IF EXISTS chat;
DROP TABLE IF EXISTS "user";
-- +goose StatementEnd 