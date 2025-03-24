-- +goose Up
-- +goose StatementBegin
-- Создаем временную таблицу без уникального ограничения
CREATE TABLE IF NOT EXISTS temp_friend (
    id VARCHAR PRIMARY KEY,
    name VARCHAR,
    birthday VARCHAR,
    chat_id VARCHAR,
    user_id VARCHAR,
    notify_at VARCHAR,
    created_at VARCHAR,
    updated_at VARCHAR
);

-- Копируем данные из существующей таблицы
INSERT INTO temp_friend
    SELECT * FROM friend;

-- Удаляем старую таблицу
DROP TABLE friend;

-- Переименовываем временную таблицу
ALTER TABLE temp_friend RENAME TO friend;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Создаем временную таблицу с уникальным ограничением
CREATE TABLE IF NOT EXISTS temp_friend (
    id VARCHAR PRIMARY KEY,
    name VARCHAR,
    birthday VARCHAR,
    chat_id VARCHAR,
    user_id VARCHAR,
    notify_at VARCHAR,
    created_at VARCHAR,
    updated_at VARCHAR,
    
    UNIQUE(name,chat_id)
);

-- Копируем данные из существующей таблицы
INSERT INTO temp_friend
    SELECT * FROM friend;

-- Удаляем старую таблицу
DROP TABLE friend;

-- Переименовываем временную таблицу
ALTER TABLE temp_friend RENAME TO friend;
-- +goose StatementEnd 