-- +goose Up
CREATE TABLE IF NOT EXISTS user_sessions (
  id BIGSERIAL PRIMARY KEY,
  user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  device_id VARCHAR(128) NOT NULL,
  token_hash VARCHAR(128) NOT NULL UNIQUE,
  expires_at TIMESTAMPTZ NOT NULL,
  revoked_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_user_sessions_user_id ON user_sessions(user_id);

CREATE INDEX IF NOT EXISTS idx_user_sessions_device ON user_sessions(user_id, device_id);

CREATE INDEX IF NOT EXISTS idx_user_sessions_expires_at ON user_sessions(expires_at);

CREATE TRIGGER trg_user_sessions_set_updated_at BEFORE
UPDATE ON user_sessions FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

-- +goose Down
DROP TABLE IF EXISTS user_sessions;