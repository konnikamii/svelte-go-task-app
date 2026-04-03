------------------------ Authorization ------------------------
-- name: GetPermissionScopesForUser :many
SELECT DISTINCT p.scope
FROM users_roles ur
  JOIN roles_permissions rp ON rp.role_id = ur.role_id
  JOIN permissions p ON p.id = rp.permission_id
WHERE ur.user_id = $1
  AND p.resource = $2
  AND p.action = $3;

-- name: AssignRoleByNameToUser :execrows
INSERT INTO users_roles (user_id, role_id)
SELECT $1,
  r.id
FROM roles r
WHERE r.name = $2 ON CONFLICT (user_id, role_id) DO NOTHING;