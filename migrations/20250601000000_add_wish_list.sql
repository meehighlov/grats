-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS wish (
    id VARCHAR PRIMARY KEY,
    name VARCHAR,
    chat_id VARCHAR,
    user_id VARCHAR,
    link VARCHAR,
    executor_id VARCHAR,
    price VARCHAR,
    wish_list_id VARCHAR,
    created_at VARCHAR,
    updated_at VARCHAR
);

CREATE TABLE IF NOT EXISTS wish_list (
    id VARCHAR PRIMARY KEY,
    name VARCHAR,
    user_id VARCHAR,
    chat_id VARCHAR,
    created_at VARCHAR,
    updated_at VARCHAR
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS wish;
DROP TABLE IF EXISTS wish_list;
-- +goose StatementEnd