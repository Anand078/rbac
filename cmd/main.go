package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/Anand078/rbac/internal/config"
	"github.com/Anand078/rbac/internal/database"
	"github.com/Anand078/rbac/internal/handlers"
	"github.com/Anand078/rbac/internal/middleware"
	"github.com/Anand078/rbac/internal/services"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.NewConnection(cfg.DatabaseURL, cfg.SupabaseURL, cfg.SupabaseServiceKey)
	if err != nil {
		errStr := err.Error()
		fmt.Println(errStr)
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize services
	authService := services.NewAuthService(db, cfg.JWTSecret)
	rbacService := services.NewRBACService(db)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	roleHandler := handlers.NewRoleHandler(rbacService)
	permissionHandler := handlers.NewPermissionHandler(rbacService)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret, rbacService)

	// Setup router
	router := gin.Default()

	// Public routes
	api := router.Group("/api")
	{
		// Authentication endpoints
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)
	}

	// Protected routes
	protected := api.Group("/")
	protected.Use(authMiddleware.Authenticate())
	{
		// Role management
		protected.POST("/roles/create", authMiddleware.RequireRole("admin"), roleHandler.CreateRole)
		protected.GET("/roles", roleHandler.GetAllRoles)
		protected.GET("/users/:userID/roles", roleHandler.GetUserRoles)
		protected.POST("/users/assign-role", authMiddleware.RequireRole("admin"), roleHandler.AssignRole)
		protected.DELETE("/users/:userID/roles/:roleID", authMiddleware.RequireRole("admin"), roleHandler.RemoveRole)

		// Permission management
		protected.POST("/permissions/create", authMiddleware.RequireRole("admin"), permissionHandler.CreatePermission)
		protected.GET("/permissions", permissionHandler.GetAllPermissions)
		protected.GET("/roles/:roleID/permissions", permissionHandler.GetRolePermissions)
		protected.POST("/permissions/grant", authMiddleware.RequireRole("admin"), permissionHandler.GrantPermission)
		protected.DELETE("/roles/:roleID/permissions/:permissionID", authMiddleware.RequireRole("admin"), permissionHandler.RevokePermission)

		// Example protected endpoints with specific permissions
		protected.GET("/courses", authMiddleware.Authorize("course", "read"), func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Course list"})
		})
		protected.POST("/courses", authMiddleware.Authorize("course", "create"), func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Course created"})
		})
		protected.GET("/grades", authMiddleware.Authorize("grades", "read"), func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Grades list"})
		})
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
