-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS wish (
    id VARCHAR PRIMARY KEY,
    name VARCHAR,
    description VARCHAR,
    chat_id VARCHAR,
    user_id VARCHAR,
    link VARCHAR,
    ozon_link VARCHAR,
    wb_link VARCHAR,
    locked VARCHAR,
    price VARCHAR,
    created_at VARCHAR,
    updated_at VARCHAR
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS wish;
-- +goose StatementEnd