-- +goose Up
ALTER TABLE friend ADD COLUMN city VARCHAR;

-- +goose Down
CREATE TABLE IF NOT EXISTS temp_friend (
    id VARCHAR PRIMARY KEY,
    name VARCHAR,
    birthday VARCHAR,
    chatid VARCHAR,
    userid VARCHAR,
    notifyat VARCHAR,
    createdat VARCHAR,
    updatedat VARCHAR,
    UNIQUE(name, chatid)
);
INSERT INTO temp_friend (id, name, birthday, chatid, userid, notifyat, createdat, updatedat)
SELECT id, name, birthday, chatid, userid, notifyat, createdat, updatedat FROM friend;
DROP TABLE friend;
ALTER TABLE temp_friend RENAME TO friend;