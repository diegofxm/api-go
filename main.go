package main

import (
	"log"
	"os"

	"go-api-orm/config"
	"go-api-orm/routes"
	"go-api-orm/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Asegurar que el archivo .env existe y tiene una clave JWT
	if err := services.EnsureEnvFile(); err != nil {
		log.Printf("Error ensuring .env file: %v", err)
	}

	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	// Configurar el modo de Gin
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Inicializar la base de datos
	config.InitDB()

	// Inicializar el router
	r := gin.Default()

	// Configurar rutas
	routes.SetupUserRoutes(r)
	routes.SetupPostRoutes(r)
	routes.SetupRoleRoutes(r)

	// Iniciar el servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
