package router

import (
	"github.com/gin-gonic/gin"
	_ "github.com/jordisetiawan/insurance-auth-service/docs"
	"github.com/jordisetiawan/insurance-auth-service/internal/handler"
	"github.com/jordisetiawan/insurance-auth-service/internal/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(authHandler *handler.AuthHandler, secret string) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.ErrorHandler())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)

			// Protected routes
			protected := auth.Group("/")
			protected.Use(middleware.AuthMiddleware(secret))
			{
				protected.POST("/logout", authHandler.Logout)
				protected.GET("/me", authHandler.GetMe)
			}
		}

		// Admin Specific Routes
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(secret), middleware.RoleMiddleware("ADMIN"), middleware.PermissionMiddleware("view_admin_dashboard"))
		{
			admin.GET("/dashboard", authHandler.AdminDashboard)
		}

		// User Specific Routes
		user := api.Group("/user")
		user.Use(middleware.AuthMiddleware(secret), middleware.RoleMiddleware("USER", "ADMIN")) // Gunakan kode peran
		{
			user.GET("/dashboard", authHandler.UserDashboard)
		}

		// User Management (Admin Only)
		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware(secret), middleware.RoleMiddleware("ADMIN"), middleware.PermissionMiddleware("manage_users"))
		{
			users.GET("", authHandler.GetUsers)
			users.GET("/:id", authHandler.GetUserByID)
			users.PUT("/:id", authHandler.UpdateUser)
			users.PATCH("/:id/role", authHandler.AssignRole)
		}

		// Role & Permission Management
		roles := api.Group("/roles")
		roles.Use(middleware.AuthMiddleware(secret), middleware.RoleMiddleware("ADMIN"), middleware.PermissionMiddleware("manage_roles"))
		{
			roles.POST("/:code/permissions", authHandler.AssignPermission)
		}
	}

	return r
}
