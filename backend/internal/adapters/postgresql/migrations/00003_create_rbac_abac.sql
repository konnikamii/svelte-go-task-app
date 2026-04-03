-- +goose Up
CREATE TABLE IF NOT EXISTS roles (
  id SMALLSERIAL PRIMARY KEY,
  name VARCHAR(64) NOT NULL UNIQUE,
  description TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS permissions (
  id BIGSERIAL PRIMARY KEY,
  resource VARCHAR(64) NOT NULL,
  ACTION VARCHAR(16) NOT NULL,
  scope VARCHAR(16) NOT NULL,
  description TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT permissions_action_check CHECK (ACTION IN ('read', 'write')),
  CONSTRAINT permissions_scope_check CHECK (scope IN ('any', 'own')),
  CONSTRAINT permissions_unique_resource_action_scope UNIQUE (resource, ACTION, scope)
);

CREATE TABLE IF NOT EXISTS users_roles (
  user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  role_id SMALLINT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (user_id, role_id)
);

CREATE TABLE IF NOT EXISTS roles_permissions (
  role_id SMALLINT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
  permission_id BIGINT NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (role_id, permission_id)
);

CREATE TRIGGER trg_roles_set_updated_at BEFORE
UPDATE ON roles FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

CREATE TRIGGER trg_permissions_set_updated_at BEFORE
UPDATE ON permissions FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

CREATE TRIGGER trg_users_roles_set_updated_at BEFORE
UPDATE ON users_roles FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

CREATE TRIGGER trg_roles_permissions_set_updated_at BEFORE
UPDATE ON roles_permissions FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

INSERT INTO roles (name, description)
VALUES ('admin', 'Full access to all resources'),
  ('user', 'Can read and write only own resources') ON CONFLICT (name) DO
UPDATE
SET description = EXCLUDED.description,
  updated_at = CURRENT_TIMESTAMP;

INSERT INTO permissions (resource, ACTION, scope, description)
VALUES ('user', 'read', 'any', 'Read any user'),
  ('user', 'write', 'any', 'Write any user'),
  ('task', 'read', 'any', 'Read any task'),
  ('task', 'write', 'any', 'Write any task'),
  ('user', 'read', 'own', 'Read own user'),
  ('user', 'write', 'own', 'Write own user'),
  ('task', 'read', 'own', 'Read own task'),
  ('task', 'write', 'own', 'Write own task') ON CONFLICT (resource, ACTION, scope) DO
UPDATE
SET description = EXCLUDED.description,
  updated_at = CURRENT_TIMESTAMP;

INSERT INTO roles_permissions (role_id, permission_id)
SELECT r.id,
  p.id
FROM roles r
  JOIN permissions p ON (
    (
      r.name = 'admin'
      AND p.scope = 'any'
    )
    OR (
      r.name = 'user'
      AND p.scope = 'own'
    )
  ) ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO users_roles (user_id, role_id)
SELECT u.id,
  r.id
FROM users u
  JOIN roles r ON r.name = 'user' ON CONFLICT (user_id, role_id) DO NOTHING;

-- +goose Down
DROP TABLE IF EXISTS roles_permissions;

DROP TABLE IF EXISTS users_roles;

DROP TABLE IF EXISTS permissions;

DROP TABLE IF EXISTS roles;