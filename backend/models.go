package main

type Permission struct {
	ID   string `json:"id"`
	Name string `json:"name"` // e.g. "user:read", "user:write"
}

type Role struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
	Policies    []string `json:"policies"`
}

type User struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Password string   `json:"password,omitempty"`
	Roles    []string `json:"roles"` // role IDs
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AssignRolesRequest struct {
	RoleIDs []string `json:"role_ids" binding:"required"`
}

type AssignPermissionsRequest struct {
	PermissionIDs []string `json:"permission_ids" binding:"required"`
}

type Policy struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Type  string `json:"type"`  // "horario" | "ip"
	Value string `json:"value"` // "08:00-18:00" | "127.0.0.1,::1"
}

type AssignPoliciesRequest struct {
	PolicyIDs []string `json:"policy_ids" binding:"required"`
}
