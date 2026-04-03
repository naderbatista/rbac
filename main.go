package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	store := NewStore()
	h := NewHandlers(store)

	r := gin.Default()

	r.POST("/login", h.Login)

	api := r.Group("/api", Auth())
	{
		users := api.Group("/users")
		users.GET("", RequirePermission(store, "user:read"), h.ListUsers)
		users.POST("", RequirePermission(store, "user:write"), h.CreateUser)
		users.PUT("/:id/roles", RequirePermission(store, "role:write"), h.AssignRoles)

		roles := api.Group("/roles")
		roles.GET("", RequirePermission(store, "role:read"), h.ListRoles)
		roles.POST("", RequirePermission(store, "role:write"), h.CreateRole)
		roles.PUT("/:id/permissions", RequirePermission(store, "role:write"), h.AssignPermissions)

		perms := api.Group("/permissions")
		perms.GET("", RequirePermission(store, "role:read"), h.ListPermissions)
		perms.POST("", RequirePermission(store, "role:write"), h.CreatePermission)
	}

	log.Println("starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
