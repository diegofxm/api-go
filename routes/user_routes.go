package routes

import (
	"github.com/gin-gonic/gin"
	"go-api-orm/controllers"
	"go-api-orm/middleware"
)

func SetupUserRoutes(router *gin.Engine) {
	api := router.Group("/api")

	// Rutas públicas
	api.POST("/register", controllers.Register)
	api.POST("/login", controllers.Login)

	// Rutas protegidas
	protected := api.Group("/users")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("", controllers.GetUsers)
		protected.GET("/:id", controllers.GetUser)
		protected.PUT("/:id", controllers.UpdateUser)
		protected.DELETE("/:id", middleware.RoleMiddleware("admin"), controllers.DeleteUser)
	}
}
