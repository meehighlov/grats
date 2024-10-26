-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS chat (
    id VARCHAR PRIMARY KEY,
    chatid VARCHAR,
    chattype VARCHAR,
    botinvitedbyid INTEGER,  -- telegram user id
    createdat VARCHAR,
	updatedat VARCHAR,

    UNIQUE(chatid)
    );

INSERT INTO chat (id, chatid, chattype, botinvitedbyid, createdat, updatedat)
SELECT RANDOM(), chatid, 'private', tgid, createdat, updatedat FROM user;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS chat;
-- +goose StatementEnd
