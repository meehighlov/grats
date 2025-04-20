-- +goose Up
-- +goose StatementBegin
-- Обновляем таблицу user
ALTER TABLE "user" 
    ALTER COLUMN created_at TYPE TIMESTAMP WITH TIME ZONE USING CASE 
        WHEN created_at ~ '^\d{2}\.\d{2}\.\d{4}T\d{2}:\d{2}:\d{2}$' THEN 
            to_timestamp(created_at, 'DD.MM.YYYY"T"HH24:MI:SS') AT TIME ZONE 'UTC'
        ELSE 
            created_at::TIMESTAMP WITH TIME ZONE 
    END,
    ALTER COLUMN updated_at TYPE TIMESTAMP WITH TIME ZONE USING CASE 
        WHEN updated_at ~ '^\d{2}\.\d{2}\.\d{4}T\d{2}:\d{2}:\d{2}$' THEN 
            to_timestamp(updated_at, 'DD.MM.YYYY"T"HH24:MI:SS') AT TIME ZONE 'UTC'
        ELSE 
            updated_at::TIMESTAMP WITH TIME ZONE 
    END;

-- Обновляем таблицу chat
ALTER TABLE "chat" 
    ALTER COLUMN created_at TYPE TIMESTAMP WITH TIME ZONE USING CASE 
        WHEN created_at ~ '^\d{2}\.\d{2}\.\d{4}T\d{2}:\d{2}:\d{2}$' THEN 
            to_timestamp(created_at, 'DD.MM.YYYY"T"HH24:MI:SS') AT TIME ZONE 'UTC'
        ELSE 
            created_at::TIMESTAMP WITH TIME ZONE 
    END,
    ALTER COLUMN updated_at TYPE TIMESTAMP WITH TIME ZONE USING CASE 
        WHEN updated_at ~ '^\d{2}\.\d{2}\.\d{4}T\d{2}:\d{2}:\d{2}$' THEN 
            to_timestamp(updated_at, 'DD.MM.YYYY"T"HH24:MI:SS') AT TIME ZONE 'UTC'
        ELSE 
            updated_at::TIMESTAMP WITH TIME ZONE 
    END;

-- Обновляем таблицу friend
ALTER TABLE "friend" 
    ALTER COLUMN created_at TYPE TIMESTAMP WITH TIME ZONE USING CASE 
        WHEN created_at ~ '^\d{2}\.\d{2}\.\d{4}T\d{2}:\d{2}:\d{2}$' THEN 
            to_timestamp(created_at, 'DD.MM.YYYY"T"HH24:MI:SS') AT TIME ZONE 'UTC'
        ELSE 
            created_at::TIMESTAMP WITH TIME ZONE 
    END,
    ALTER COLUMN updated_at TYPE TIMESTAMP WITH TIME ZONE USING CASE 
        WHEN updated_at ~ '^\d{2}\.\d{2}\.\d{4}T\d{2}:\d{2}:\d{2}$' THEN 
            to_timestamp(updated_at, 'DD.MM.YYYY"T"HH24:MI:SS') AT TIME ZONE 'UTC'
        ELSE 
            updated_at::TIMESTAMP WITH TIME ZONE 
    END;

-- Обновляем таблицу wish
ALTER TABLE "wish" 
    ALTER COLUMN created_at TYPE TIMESTAMP WITH TIME ZONE USING CASE 
        WHEN created_at ~ '^\d{2}\.\d{2}\.\d{4}T\d{2}:\d{2}:\d{2}$' THEN 
            to_timestamp(created_at, 'DD.MM.YYYY"T"HH24:MI:SS') AT TIME ZONE 'UTC'
        ELSE 
            created_at::TIMESTAMP WITH TIME ZONE 
    END,
    ALTER COLUMN updated_at TYPE TIMESTAMP WITH TIME ZONE USING CASE 
        WHEN updated_at ~ '^\d{2}\.\d{2}\.\d{4}T\d{2}:\d{2}:\d{2}$' THEN 
            to_timestamp(updated_at, 'DD.MM.YYYY"T"HH24:MI:SS') AT TIME ZONE 'UTC'
        ELSE 
            updated_at::TIMESTAMP WITH TIME ZONE 
    END;

-- Обновляем таблицу wish_list
ALTER TABLE "wish_list" 
    ALTER COLUMN created_at TYPE TIMESTAMP WITH TIME ZONE USING CASE 
        WHEN created_at ~ '^\d{2}\.\d{2}\.\d{4}T\d{2}:\d{2}:\d{2}$' THEN 
            to_timestamp(created_at, 'DD.MM.YYYY"T"HH24:MI:SS') AT TIME ZONE 'UTC'
        ELSE 
            created_at::TIMESTAMP WITH TIME ZONE 
    END,
    ALTER COLUMN updated_at TYPE TIMESTAMP WITH TIME ZONE USING CASE 
        WHEN updated_at ~ '^\d{2}\.\d{2}\.\d{4}T\d{2}:\d{2}:\d{2}$' THEN 
            to_timestamp(updated_at, 'DD.MM.YYYY"T"HH24:MI:SS') AT TIME ZONE 'UTC'
        ELSE 
            updated_at::TIMESTAMP WITH TIME ZONE 
    END;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Возвращаем таблицу user к предыдущему состоянию
ALTER TABLE "user" 
    ALTER COLUMN created_at TYPE VARCHAR(255),
    ALTER COLUMN updated_at TYPE VARCHAR(255);

-- Возвращаем таблицу chat к предыдущему состоянию
ALTER TABLE "chat" 
    ALTER COLUMN created_at TYPE VARCHAR(255),
    ALTER COLUMN updated_at TYPE VARCHAR(255);

-- Возвращаем таблицу friend к предыдущему состоянию
ALTER TABLE "friend" 
    ALTER COLUMN created_at TYPE VARCHAR(255),
    ALTER COLUMN updated_at TYPE VARCHAR(255);

-- Возвращаем таблицу wish к предыдущему состоянию
ALTER TABLE "wish" 
    ALTER COLUMN created_at TYPE VARCHAR(255),
    ALTER COLUMN updated_at TYPE VARCHAR(255);

-- Возвращаем таблицу wish_list к предыдущему состоянию
ALTER TABLE "wish_list" 
    ALTER COLUMN created_at TYPE VARCHAR(255),
    ALTER COLUMN updated_at TYPE VARCHAR(255);
-- +goose StatementEnd 