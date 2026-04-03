package main

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type Store struct {
	mu          sync.RWMutex
	users       map[string]User
	roles       map[string]Role
	permissions map[string]Permission
}

func NewStore() *Store {
	s := &Store{
		users:       make(map[string]User),
		roles:       make(map[string]Role),
		permissions: make(map[string]Permission),
	}
	s.seed()
	return s
}

func (s *Store) seed() {
	perms := []Permission{
		{ID: uuid.NewString(), Name: "user:read"},
		{ID: uuid.NewString(), Name: "user:write"},
		{ID: uuid.NewString(), Name: "role:read"},
		{ID: uuid.NewString(), Name: "role:write"},
	}
	for _, p := range perms {
		s.permissions[p.ID] = p
	}

	adminRole := Role{
		ID:   uuid.NewString(),
		Name: "admin",
	}
	for _, p := range perms {
		adminRole.Permissions = append(adminRole.Permissions, p.ID)
	}
	s.roles[adminRole.ID] = adminRole

	viewerRole := Role{
		ID:          uuid.NewString(),
		Name:        "viewer",
		Permissions: []string{perms[0].ID, perms[2].ID}, // user:read, role:read
	}
	s.roles[viewerRole.ID] = viewerRole

	admin := User{
		ID:       uuid.NewString(),
		Username: "admin",
		Password: "admin123",
		Roles:    []string{adminRole.ID},
	}
	s.users[admin.ID] = admin
}

// --- Users ---

func (s *Store) CreateUser(username, password string) User {
	s.mu.Lock()
	defer s.mu.Unlock()

	u := User{ID: uuid.NewString(), Username: username, Password: password}
	s.users[u.ID] = u
	return u
}

func (s *Store) GetUser(id string) (User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.users[id]
	return u, ok
}

func (s *Store) ListUsers() []User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]User, 0, len(s.users))
	for _, u := range s.users {
		out = append(out, u)
	}
	return out
}

func (s *Store) FindUserByUsername(username string) (User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, u := range s.users {
		if u.Username == username {
			return u, true
		}
	}
	return User{}, false
}

func (s *Store) AssignRoles(userID string, roleIDs []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	u, ok := s.users[userID]
	if !ok {
		return fmt.Errorf("user not found")
	}
	for _, rid := range roleIDs {
		if _, exists := s.roles[rid]; !exists {
			return fmt.Errorf("role %s not found", rid)
		}
	}
	u.Roles = roleIDs
	s.users[userID] = u
	return nil
}

// --- Roles ---

func (s *Store) CreateRole(name string) Role {
	s.mu.Lock()
	defer s.mu.Unlock()

	r := Role{ID: uuid.NewString(), Name: name}
	s.roles[r.ID] = r
	return r
}

func (s *Store) ListRoles() []Role {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]Role, 0, len(s.roles))
	for _, r := range s.roles {
		out = append(out, r)
	}
	return out
}

func (s *Store) AssignPermissions(roleID string, permIDs []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	r, ok := s.roles[roleID]
	if !ok {
		return fmt.Errorf("role not found")
	}
	for _, pid := range permIDs {
		if _, exists := s.permissions[pid]; !exists {
			return fmt.Errorf("permission %s not found", pid)
		}
	}
	r.Permissions = permIDs
	s.roles[roleID] = r
	return nil
}

// --- Permissions ---

func (s *Store) CreatePermission(name string) Permission {
	s.mu.Lock()
	defer s.mu.Unlock()

	p := Permission{ID: uuid.NewString(), Name: name}
	s.permissions[p.ID] = p
	return p
}

func (s *Store) ListPermissions() []Permission {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]Permission, 0, len(s.permissions))
	for _, p := range s.permissions {
		out = append(out, p)
	}
	return out
}

// UserHasPermission resolves the full chain: user -> roles -> permissions.
func (s *Store) UserHasPermission(userID, permName string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	u, ok := s.users[userID]
	if !ok {
		return false
	}
	for _, rid := range u.Roles {
		role, ok := s.roles[rid]
		if !ok {
			continue
		}
		for _, pid := range role.Permissions {
			if perm, ok := s.permissions[pid]; ok && perm.Name == permName {
				return true
			}
		}
	}
	return false
}
