package routes

import (
	"github.com/gin-gonic/gin"
	"go-api-orm/controllers"
	"go-api-orm/middleware"
)

func SetupRoleRoutes(router *gin.Engine) {
	api := router.Group("/api")

	// Rutas de roles (todas protegidas ya que solo los administradores deberían manejar roles)
	roles := api.Group("/roles")
	roles.Use(middleware.AuthMiddleware())
	{
		roles.POST("", controllers.CreateRole)
		roles.GET("", controllers.GetRoles)
		roles.GET("/:id", controllers.GetRole)
		roles.PUT("/:id", controllers.UpdateRole)
		roles.DELETE("/:id", controllers.DeleteRole)
	}
}
