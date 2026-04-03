package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	store := NewStore()
	h := NewHandlers(store)

	r := gin.Default()
	r.Use(CORS())

	r.POST("/login", h.Login)

	api := r.Group("/api", Auth())
	{
		api.GET("/me", h.Me)

		// ABAC applies to all data routes below
		protected := api.Group("", RequireABAC(store))

		users := protected.Group("/users")
		users.GET("", RequirePermission(store, "user:read"), h.ListUsers)
		users.POST("", RequirePermission(store, "user:write"), h.CreateUser)
		users.PUT("/:id/roles", RequirePermission(store, "role:write"), h.AssignRoles)

		roles := protected.Group("/roles")
		roles.GET("", RequirePermission(store, "role:read"), h.ListRoles)
		roles.POST("", RequirePermission(store, "role:write"), h.CreateRole)
		roles.PUT("/:id/permissions", RequirePermission(store, "role:write"), h.AssignPermissions)
		roles.PUT("/:id/policies", RequirePermission(store, "role:write"), h.AssignPolicies)

		perms := protected.Group("/permissions")
		perms.GET("", RequirePermission(store, "role:read"), h.ListPermissions)
		perms.POST("", RequirePermission(store, "role:write"), h.CreatePermission)

		policies := protected.Group("/policies")
		policies.GET("", RequirePermission(store, "role:read"), h.ListPolicies)
		policies.POST("", RequirePermission(store, "role:write"), h.CreatePolicyHandler)
	}

	log.Println("starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
