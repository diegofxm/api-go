package routes

import (
	"github.com/gin-gonic/gin"
	"go-api-orm/controllers"
	"go-api-orm/middleware"
)

func SetupPostRoutes(router *gin.Engine) {
	api := router.Group("/api")

	// Rutas públicas de posts
	posts := api.Group("/posts")
	{
		posts.GET("", controllers.GetPosts)
		posts.GET("/:slug", controllers.GetPostBySlug)

		// Rutas protegidas que requieren autenticación
		protected := posts.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.POST("", controllers.CreatePost)
			protected.PUT("/:slug", controllers.UpdatePost)
			protected.DELETE("/:slug", controllers.DeletePost)
		}
	}
}
