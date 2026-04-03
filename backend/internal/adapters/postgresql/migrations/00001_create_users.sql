-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_modified_column() RETURNS TRIGGER AS $$ BEGIN NEW.updated_at = NOW();

RETURN NEW;

END;

$$ LANGUAGE plpgsql;

-- +goose StatementEnd
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    PASSWORD VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER update_modified_time BEFORE
UPDATE ON users FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

-- +goose Down
DROP TABLE IF EXISTS users;

DROP FUNCTION IF EXISTS update_modified_column;