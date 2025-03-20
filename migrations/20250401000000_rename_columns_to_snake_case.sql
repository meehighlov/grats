-- +goose Up
-- +goose StatementBegin
-- Переименование колонок в таблице user
ALTER TABLE user RENAME COLUMN tgid TO tg_id;
ALTER TABLE user RENAME COLUMN tgusername TO tg_username;
ALTER TABLE user RENAME COLUMN chatid TO chat_id;
ALTER TABLE user RENAME COLUMN isadmin TO is_admin;
ALTER TABLE user RENAME COLUMN createdat TO created_at;
ALTER TABLE user RENAME COLUMN updatedat TO updated_at;

-- Переименование колонок в таблице friend
ALTER TABLE friend RENAME COLUMN chatid TO chat_id;
ALTER TABLE friend RENAME COLUMN userid TO user_id;
ALTER TABLE friend RENAME COLUMN notifyat TO notify_at;
ALTER TABLE friend RENAME COLUMN createdat TO created_at;
ALTER TABLE friend RENAME COLUMN updatedat TO updated_at;

-- Переименование колонок в таблице chat
ALTER TABLE chat RENAME COLUMN chatid TO chat_id;
ALTER TABLE chat RENAME COLUMN chattype TO chat_type;
ALTER TABLE chat RENAME COLUMN botinvitedbyid TO bot_invited_by_id;
ALTER TABLE chat RENAME COLUMN createdat TO created_at;
ALTER TABLE chat RENAME COLUMN updatedat TO updated_at;

-- Переименование колонок в таблице access
ALTER TABLE access RENAME COLUMN tgusername TO tg_username;
ALTER TABLE access RENAME COLUMN createdat TO created_at;
ALTER TABLE access RENAME COLUMN updatedat TO updated_at;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Возврат на исходное состояние для таблицы user
ALTER TABLE user RENAME COLUMN tg_id TO tgid;
ALTER TABLE user RENAME COLUMN tg_username TO tgusername;
ALTER TABLE user RENAME COLUMN chat_id TO chatid;
ALTER TABLE user RENAME COLUMN is_admin TO isadmin;
ALTER TABLE user RENAME COLUMN created_at TO createdat;
ALTER TABLE user RENAME COLUMN updated_at TO updatedat;

-- Возврат на исходное состояние для таблицы friend
ALTER TABLE friend RENAME COLUMN chat_id TO chatid;
ALTER TABLE friend RENAME COLUMN user_id TO userid;
ALTER TABLE friend RENAME COLUMN notify_at TO notifyat;
ALTER TABLE friend RENAME COLUMN created_at TO createdat;
ALTER TABLE friend RENAME COLUMN updated_at TO updatedat;

-- Возврат на исходное состояние для таблицы chat
ALTER TABLE chat RENAME COLUMN chat_id TO chatid;
ALTER TABLE chat RENAME COLUMN chat_type TO chattype;
ALTER TABLE chat RENAME COLUMN bot_invited_by_id TO botinvitedbyid;
ALTER TABLE chat RENAME COLUMN created_at TO createdat;
ALTER TABLE chat RENAME COLUMN updated_at TO updatedat;

-- Возврат на исходное состояние для таблицы access
ALTER TABLE access RENAME COLUMN tg_username TO tgusername;
ALTER TABLE access RENAME COLUMN created_at TO createdat;
ALTER TABLE access RENAME COLUMN updated_at TO updatedat;
-- +goose StatementEnd