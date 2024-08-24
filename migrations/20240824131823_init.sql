-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user (
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
CREATE TABLE IF NOT EXISTS friend (
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
CREATE TABLE IF NOT EXISTS access (
		id VARCHAR PRIMARY KEY,
		tgusername VARCHAR,
		createdat VARCHAR,
		updatedat VARCHAR,

		UNIQUE(tgusername)
	);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user;
DROP TABLE IF EXISTS friend;
DROP TABLE IF EXISTS access;
-- +goose StatementEnd
