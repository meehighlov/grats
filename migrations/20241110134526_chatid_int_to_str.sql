-- +goose Up
-- +goose StatementBegin

-------------------------------------------- FRIEND --------------------------------------------

CREATE TABLE IF NOT EXISTS temp_friend (
		id VARCHAR PRIMARY KEY,
		name VARCHAR,
		birthday VARCHAR,
		chatid VARCHAR,
		userid VARCHAR,
		notifyat VARCHAR,
		createdat VARCHAR,
		updatedat VARCHAR,

		UNIQUE(name,chatid)
	);

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
    CAST(chatid as TEXT),
    CAST(userid as TEXT),
    notifyat,
    createdat,
    updatedat
FROM friend;

DROP TABLE friend;
ALTER TABLE temp_friend RENAME TO friend;

-------------------------------------------- USER --------------------------------------------

CREATE TABLE IF NOT EXISTS temp_user (
		id VARCHAR PRIMARY KEY,
		tgid VARCHAR,
		name VARCHAR,
		tgusername VARCHAR,
		chatid VARCHAR,
		birthday VARCHAR,
		isadmin INTEGER,
		createdat VARCHAR,
		updatedat VARCHAR,

		UNIQUE(tgid)
	);

INSERT INTO temp_user (
    id,
    tgid,
    name,
    tgusername,
    chatid,
    birthday,
    isadmin,
    createdat,
    updatedat
)
SELECT
    id,
    CAST(tgid as TEXT),
    name,
    tgusername,
    chatid,
    birthday,
    isadmin,
    createdat,
    updatedat
FROM user;

DROP TABLE user;
ALTER TABLE temp_user RENAME TO user;

-------------------------------------------- CHAT --------------------------------------------

CREATE TABLE IF NOT EXISTS temp_chat (
    id VARCHAR PRIMARY KEY,
    chatid VARCHAR,
    chattype VARCHAR,
    botinvitedbyid VARCHAR,  -- telegram user id
    createdat VARCHAR,
	updatedat VARCHAR,

    UNIQUE(chatid)
    );

INSERT INTO temp_chat (
    id,
    chatid,
    chattype,
    botinvitedbyid,  -- telegram user id
    createdat,
	updatedat
)
SELECT
    id,
    chatid,
    chattype,
    CAST(botinvitedbyid as TEXT),  -- telegram user id
    createdat,
	updatedat
FROM chat;

DROP table chat;
ALTER TABLE temp_chat RENAME TO chat;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-------------------------------------------- FRIEND --------------------------------------------

CREATE TABLE IF NOT EXISTS temp_friend (
		id VARCHAR PRIMARY KEY,
		name VARCHAR,
		birthday VARCHAR,
		chatid INTEGER,
		userid INTEGER,
		notifyat VARCHAR,
		createdat VARCHAR,
		updatedat VARCHAR,

		UNIQUE(name,chatid)
	);

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
    CAST(chatid as INTEGER),
    CAST(userid as INTEGER),
    notifyat,
    createdat,
    updatedat
FROM friend;

DROP TABLE friend;
ALTER TABLE temp_friend RENAME TO friend;

-------------------------------------------- USER --------------------------------------------

CREATE TABLE IF NOT EXISTS temp_user (
		id VARCHAR PRIMARY KEY,
		tgid INTEGER,
		name VARCHAR,
		tgusername VARCHAR,
		chatid VARCHAR,
		birthday VARCHAR,
		isadmin INTEGER,
		createdat VARCHAR,
		updatedat VARCHAR,

		UNIQUE(tgid)
	);

INSERT INTO temp_user (
    id,
    tgid,
    name,
    tgusername,
    chatid,
    birthday,
    isadmin,
    createdat,
    updatedat
)
SELECT
    id,
    CAST(tgid as INTEGER),
    name,
    tgusername,
    chatid,
    birthday,
    isadmin,
    createdat,
    updatedat
FROM user;

DROP TABLE user;
ALTER TABLE temp_user RENAME TO user;

-------------------------------------------- CHAT --------------------------------------------

CREATE TABLE IF NOT EXISTS temp_chat (
    id VARCHAR PRIMARY KEY,
    chatid VARCHAR,
    chattype VARCHAR,
    botinvitedbyid INTEGER,  -- telegram user id
    createdat VARCHAR,
	updatedat VARCHAR,

    UNIQUE(chatid)
    );

INSERT INTO temp_chat (
    id,
    chatid,
    chattype,
    botinvitedbyid,  -- telegram user id
    createdat,
	updatedat
)
SELECT
    id,
    chatid,
    chattype,
    CAST(botinvitedbyid as INTEGER),  -- telegram user id
    createdat,
	updatedat
FROM chat;

DROP table chat;
ALTER TABLE temp_chat RENAME TO chat;

-- +goose StatementEnd
