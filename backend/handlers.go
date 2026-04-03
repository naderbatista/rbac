package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	store *Store
}

func NewHandlers(s *Store) *Handlers {
	return &Handlers{store: s}
}

// --- Auth ---

func (h *Handlers) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, ok := h.store.FindUserByUsername(req.Username)
	if !ok || user.Password != req.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "credenciais inválidas"})
		return
	}

	token, err := GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao gerar token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *Handlers) Me(c *gin.Context) {
	userID := c.GetString("userID")
	user, ok := h.store.GetUser(userID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "usuário não encontrado"})
		return
	}
	user.Password = ""
	c.JSON(http.StatusOK, gin.H{
		"user":        user,
		"permissions": h.store.UserPermissionNames(userID),
		"policies":    h.store.UserPolicies(userID),
	})
}

// --- Users ---

func (h *Handlers) ListUsers(c *gin.Context) {
	users := h.store.ListUsers()
	// strip passwords from response
	for i := range users {
		users[i].Password = ""
	}
	c.JSON(http.StatusOK, users)
}

func (h *Handlers) CreateUser(c *gin.Context) {
	var body struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, exists := h.store.FindUserByUsername(body.Username); exists {
		c.JSON(http.StatusConflict, gin.H{"error": "nome de usuário já existe"})
		return
	}

	u := h.store.CreateUser(body.Username, body.Password)
	u.Password = ""
	c.JSON(http.StatusCreated, u)
}

func (h *Handlers) AssignRoles(c *gin.Context) {
	targetID := c.Param("id")
	if targetID == c.GetString("userID") {
		c.JSON(http.StatusForbidden, gin.H{"error": "não é possível modificar seus próprios perfis"})
		return
	}

	var req AssignRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.store.AssignRoles(targetID, req.RoleIDs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "perfis atribuídos"})
}

// --- Roles ---

func (h *Handlers) ListRoles(c *gin.Context) {
	c.JSON(http.StatusOK, h.store.ListRoles())
}

func (h *Handlers) CreateRole(c *gin.Context) {
	var body struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, h.store.CreateRole(body.Name))
}

func (h *Handlers) AssignPermissions(c *gin.Context) {
	roleID := c.Param("id")
	var req AssignPermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.PermissionIDs) == 0 {
		if role, ok := h.store.GetRole(roleID); ok && role.Name == "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "não é possível remover todas as permissões do perfil admin"})
			return
		}
	}

	if err := h.store.AssignPermissions(roleID, req.PermissionIDs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "permissões atribuídas"})
}

// --- Permissions ---

func (h *Handlers) ListPermissions(c *gin.Context) {
	c.JSON(http.StatusOK, h.store.ListPermissions())
}

func (h *Handlers) CreatePermission(c *gin.Context) {
	var body struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, h.store.CreatePermission(body.Name))
}

// --- Policies (ABAC) ---

func (h *Handlers) ListPolicies(c *gin.Context) {
	c.JSON(http.StatusOK, h.store.ListPolicies())
}

func (h *Handlers) CreatePolicyHandler(c *gin.Context) {
	var body struct {
		Name  string `json:"name" binding:"required"`
		Type  string `json:"type" binding:"required"`
		Value string `json:"value" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if body.Type != "horario" && body.Type != "ip" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tipo deve ser 'horario' ou 'ip'"})
		return
	}
	c.JSON(http.StatusCreated, h.store.CreatePolicy(body.Name, body.Type, body.Value))
}

func (h *Handlers) AssignPolicies(c *gin.Context) {
	var req AssignPoliciesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.store.AssignPolicies(c.Param("id"), req.PolicyIDs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "políticas atribuídas"})
}
