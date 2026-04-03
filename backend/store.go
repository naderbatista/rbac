package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Store struct {
	mu          sync.RWMutex
	users       map[string]User
	roles       map[string]Role
	permissions map[string]Permission
	policies    map[string]Policy
}

func NewStore() *Store {
	s := &Store{
		users:       make(map[string]User),
		roles:       make(map[string]Role),
		permissions: make(map[string]Permission),
		policies:    make(map[string]Policy),
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

	// ABAC policies
	horario := Policy{ID: uuid.NewString(), Name: "horário comercial", Type: "horario", Value: "08:00-18:00"}
	redeLocal := Policy{ID: uuid.NewString(), Name: "rede interna", Type: "ip", Value: "127.0.0.1,::1"}
	s.policies[horario.ID] = horario
	s.policies[redeLocal.ID] = redeLocal

	// Viewer restrito por ABAC; admin sem restrições
	viewerRole.Policies = []string{horario.ID, redeLocal.ID}
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
		return fmt.Errorf("usuário não encontrado")
	}
	for _, rid := range roleIDs {
		if _, exists := s.roles[rid]; !exists {
			return fmt.Errorf("perfil %s não encontrado", rid)
		}
	}
	u.Roles = roleIDs
	s.users[userID] = u
	return nil
}

// --- Roles ---

func (s *Store) GetRole(id string) (Role, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.roles[id]
	return r, ok
}

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
		return fmt.Errorf("perfil não encontrado")
	}
	for _, pid := range permIDs {
		if _, exists := s.permissions[pid]; !exists {
			return fmt.Errorf("permissão %s não encontrada", pid)
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

// UserPermissionNames returns all resolved permission names for a user.
func (s *Store) UserPermissionNames(userID string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	u, ok := s.users[userID]
	if !ok {
		return nil
	}
	seen := map[string]bool{}
	var out []string
	for _, rid := range u.Roles {
		role, ok := s.roles[rid]
		if !ok {
			continue
		}
		for _, pid := range role.Permissions {
			if perm, ok := s.permissions[pid]; ok && !seen[perm.Name] {
				seen[perm.Name] = true
				out = append(out, perm.Name)
			}
		}
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

// --- Policies ---

func (s *Store) CreatePolicy(name, ptype, value string) Policy {
	s.mu.Lock()
	defer s.mu.Unlock()

	p := Policy{ID: uuid.NewString(), Name: name, Type: ptype, Value: value}
	s.policies[p.ID] = p
	return p
}

func (s *Store) ListPolicies() []Policy {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]Policy, 0, len(s.policies))
	for _, p := range s.policies {
		out = append(out, p)
	}
	return out
}

func (s *Store) AssignPolicies(roleID string, policyIDs []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	r, ok := s.roles[roleID]
	if !ok {
		return fmt.Errorf("perfil não encontrado")
	}
	for _, pid := range policyIDs {
		if _, exists := s.policies[pid]; !exists {
			return fmt.Errorf("política %s não encontrada", pid)
		}
	}
	r.Policies = policyIDs
	s.roles[roleID] = r
	return nil
}

// UserPolicies returns all policies applied to a user through their roles.
func (s *Store) UserPolicies(userID string) []Policy {
	s.mu.RLock()
	defer s.mu.RUnlock()

	u, ok := s.users[userID]
	if !ok {
		return nil
	}
	seen := map[string]bool{}
	var out []Policy
	for _, rid := range u.Roles {
		role, ok := s.roles[rid]
		if !ok {
			continue
		}
		for _, pid := range role.Policies {
			if p, ok := s.policies[pid]; ok && !seen[p.ID] {
				seen[p.ID] = true
				out = append(out, p)
			}
		}
	}
	return out
}

// EvaluateABAC checks all policies attached to the user's roles.
// Returns (true, "") if all pass, or (false, reason) on first failure.
func (s *Store) EvaluateABAC(userID, clientIP string, now time.Time) (bool, string) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	u, ok := s.users[userID]
	if !ok {
		return false, "usuário não encontrado"
	}
	for _, rid := range u.Roles {
		role, ok := s.roles[rid]
		if !ok {
			continue
		}
		for _, pid := range role.Policies {
			policy, ok := s.policies[pid]
			if !ok {
				continue
			}
			if !evaluatePolicy(policy, clientIP, now) {
				return false, fmt.Sprintf("política '%s' negou acesso", policy.Name)
			}
		}
	}
	return true, ""
}
