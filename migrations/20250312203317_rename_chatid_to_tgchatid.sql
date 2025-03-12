-- +goose Up
-- +goose StatementBegin

-------------------------------------------- CHAT --------------------------------------------

-- Создаем временную таблицу с новой структурой
CREATE TABLE IF NOT EXISTS temp_chat (
    id VARCHAR PRIMARY KEY,
    tgchatid VARCHAR,  -- переименовываем chatid в tgchatid
    chattype VARCHAR,
    botinvitedbyid VARCHAR,
    greeting_template VARCHAR,
    createdat VARCHAR,
    updatedat VARCHAR,

    UNIQUE(tgchatid)
);

-- Копируем данные из старой таблицы в новую
INSERT INTO temp_chat (
    id,
    tgchatid,
    chattype,
    botinvitedbyid,
    greeting_template,
    createdat,
    updatedat
)
SELECT
    id,
    chatid,  -- старое поле chatid становится tgchatid
    chattype,
    botinvitedbyid,
    greeting_template,
    createdat,
    updatedat
FROM chat;

-- Удаляем старую таблицу и переименовываем новую
DROP TABLE chat;
ALTER TABLE temp_chat RENAME TO chat;

-------------------------------------------- FRIEND --------------------------------------------

-- Создаем временную таблицу с новой структурой
CREATE TABLE IF NOT EXISTS temp_friend (
    id VARCHAR PRIMARY KEY,
    name VARCHAR,
    birthday VARCHAR,
    chatid VARCHAR,  -- ID записи из таблицы Chat
    userid VARCHAR,
    notifyat VARCHAR,
    createdat VARCHAR,
    updatedat VARCHAR,

    UNIQUE(name,chatid)
);

-- Копируем данные из старой таблицы в новую
INSERT INTO temp_friend (
    id,
    name,
    birthday,
    chatid,
    userid,
    notifyat,
    createdat,
    updatedat
)
SELECT
    id,
    name,
    birthday,
    (SELECT c.id FROM chat c WHERE c.tgchatid = friend.chatid),  -- Устанавливаем связь с таблицей chat
    userid,
    notifyat,
    createdat,
    updatedat
FROM friend;

-- Удаляем старую таблицу и переименовываем новую
DROP TABLE friend;
ALTER TABLE temp_friend RENAME TO friend;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-------------------------------------------- FRIEND --------------------------------------------

-- Создаем временную таблицу с прежней структурой
CREATE TABLE IF NOT EXISTS temp_friend (
    id VARCHAR PRIMARY KEY,
    name VARCHAR,
    birthday VARCHAR,
    chatid VARCHAR,  -- Возвращаем к старому формату - ID чата в Telegram
    userid VARCHAR,
    notifyat VARCHAR,
    createdat VARCHAR,
    updatedat VARCHAR,

    UNIQUE(name,chatid)
);

-- Копируем данные из текущей таблицы в старый формат
INSERT INTO temp_friend (
    id,
    name,
    birthday,
    chatid,
    userid,
    notifyat,
    createdat,
    updatedat
)
SELECT
    id,
    name,
    birthday,
    (SELECT c.tgchatid FROM chat c WHERE c.id = friend.chatid),  -- Возвращаем ID чата в Telegram
    userid,
    notifyat,
    createdat,
    updatedat
FROM friend;

-- Удаляем текущую таблицу и переименовываем временную
DROP TABLE friend;
ALTER TABLE temp_friend RENAME TO friend;

-------------------------------------------- CHAT --------------------------------------------

-- Создаем временную таблицу с прежней структурой
CREATE TABLE IF NOT EXISTS temp_chat (
    id VARCHAR PRIMARY KEY,
    chatid VARCHAR,  -- Возвращаем старое имя поля
    chattype VARCHAR,
    botinvitedbyid VARCHAR,
    greeting_template VARCHAR,
    createdat VARCHAR,
    updatedat VARCHAR,

    UNIQUE(chatid)
);

-- Копируем данные из текущей таблицы в старый формат
INSERT INTO temp_chat (
    id,
    chatid,
    chattype,
    botinvitedbyid,
    greeting_template,
    createdat,
    updatedat
)
SELECT
    id,
    tgchatid,  -- tgchatid становится chatid
    chattype,
    botinvitedbyid,
    greeting_template,
    createdat,
    updatedat
FROM chat;

-- Удаляем текущую таблицу и переименовываем временную
DROP TABLE chat;
ALTER TABLE temp_chat RENAME TO chat;

-- +goose StatementEnd 