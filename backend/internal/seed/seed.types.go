package seed

type SeedTotals struct {
	Users           int64 `json:"users"`
	Tasks           int64 `json:"tasks"`
	Roles           int64 `json:"roles"`
	Permissions     int64 `json:"permissions"`
	RolePermissions int64 `json:"rolePermissions"`
	UserRoles       int64 `json:"userRoles"`
}

type SeedCredential struct {
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Roles    []string `json:"roles"`
}

type SeedResult struct {
	Totals      SeedTotals       `json:"totals"`
	Credentials []SeedCredential `json:"credentials"`
}
