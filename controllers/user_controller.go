package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go-api-orm/config"
	"go-api-orm/models"
	"go-api-orm/services"
	"go-api-orm/utils"
)

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	RoleID   uint   `json:"role" binding:"omitempty"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserDataResponse estructura específica para la respuesta de datos del usuario
type UserDataResponse struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RegisterResponse estructura para la respuesta completa
type RegisterResponse struct {
	Message string           `json:"message"`
	Data    UserDataResponse `json:"data"`
}

// Login maneja el inicio de sesión de usuarios
func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		status, response := services.ErrorResponse(services.NewAPIError(
			http.StatusBadRequest,
			"INVALID_INPUT",
			"Entrada inválida",
			err.Error(),
			nil,
		))
		c.JSON(status, response)
		return
	}

	var user models.User
	if err := config.DB.Preload("Role").Where("email = ?", input.Email).First(&user).Error; err != nil {
		status, response := services.ErrorResponse(services.NewAPIError(
			http.StatusUnauthorized,
			"INVALID_CREDENTIALS",
			"Credenciales inválidas",
			"",
			nil,
		))
		c.JSON(status, response)
		return
	}

	if err := user.CheckPassword(input.Password); err != nil {
		status, response := services.ErrorResponse(services.NewAPIError(
			http.StatusUnauthorized,
			"INVALID_CREDENTIALS",
			"Credenciales inválidas",
			"",
			nil,
		))
		c.JSON(status, response)
		return
	}

	// Generar token JWT
	token, err := utils.GenerateToken(user.ID, user.Role.Name)
	if err != nil {
		status, response := services.ErrorResponse(services.NewAPIError(
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"Error interno",
			err.Error(),
			nil,
		))
		c.JSON(status, response)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role.Name,
		},
		"token": token,
	})
}

// Register maneja el registro de nuevos usuarios
func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		status, response := services.ErrorResponse(services.NewAPIError(
			http.StatusBadRequest,
			"INVALID_INPUT",
			"Datos de entrada inválidos",
			err.Error(),
			nil,
		))
		c.JSON(status, response)
		return
	}

	// Si no se proporciona un rol, usar el rol por defecto (user)
	if input.RoleID == 0 {
		var defaultRole models.Role
		if err := config.DB.Where("name = ?", "user").First(&defaultRole).Error; err != nil {
			status, response := services.ErrorResponse(services.NewAPIError(
				http.StatusInternalServerError,
				"INTERNAL_ERROR",
				"Error interno",
				err.Error(),
				nil,
			))
			c.JSON(status, response)
			return
		}
		input.RoleID = defaultRole.ID
	} else {
		// Verificar si el rol existe
		var role models.Role
		if err := config.DB.First(&role, input.RoleID).Error; err != nil {
			status, response := services.ErrorResponse(services.NewAPIError(
				http.StatusBadRequest,
				"INVALID_ROLE",
				"Rol inválido",
				"El rol especificado no existe",
				nil,
			))
			c.JSON(status, response)
			return
		}
	}

	user := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
		RoleID:   input.RoleID,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		status, response := services.ErrorResponse(services.NewAPIError(
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"Error interno",
			err.Error(),
			nil,
		))
		c.JSON(status, response)
		return
	}

	// Recargar el usuario para obtener la relación con el rol
	config.DB.Preload("Role").First(&user, user.ID)

	response := RegisterResponse{
		Message: "Usuario creado exitosamente",
		Data: UserDataResponse{
			ID:        user.ID,
			Email:     user.Email,
			Role:      user.Role.Name,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}

	c.JSON(http.StatusCreated, response)
}

// GetUsers obtiene la lista de usuarios
func GetUsers(c *gin.Context) {
	var users []models.User
	
	// Obtener parámetros de paginación y búsqueda
	pagination := services.GeneratePaginationFromRequest(c)
	searchFilters := services.ExtractSearchParams(c)
	sortParams := services.ExtractSortParams(c)
	
	// Aplicar filtros y paginación
	db := config.DB.Preload("Role")
	db = services.ApplySearchFilters(db, searchFilters)
	db = services.ApplySorting(db, sortParams)
	
	err := db.Scopes(services.Paginate(users, &pagination, db)).Find(&users).Error
	if err != nil {
		status, response := services.ErrorResponse(services.NewAPIError(
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"Error interno",
			err.Error(),
			nil,
		))
		c.JSON(status, response)
		return
	}

	// Construir los componentes de la respuesta
	paginationResponse := services.BuildPaginationResponse(c, pagination.Page, pagination.Limit, pagination.TotalRows, pagination.TotalPages)
	
	metadataResponse := services.BuildMetadataResponse(
		[]services.SearchField{
			{
				Name:        "username",
				Type:        "string",
				Description: "Nombre de usuario",
				Operators:   []string{"eq", "like", "nlike"},
			},
			{
				Name:        "email",
				Type:        "string",
				Description: "Correo electrónico",
				Operators:   []string{"eq", "like", "nlike"},
			},
			{
				Name:        "created_at",
				Type:        "date",
				Description: "Fecha de creación",
				Operators:   []string{"gt", "gte", "lt", "lte"},
			},
		},
		[]services.SortField{
			{
				Name:        "username",
				Description: "Ordenar por nombre de usuario",
			},
			{
				Name:        "created_at",
				Description: "Ordenar por fecha de creación",
			},
		},
		searchFilters,
		sortParams,
	)

	// Construir la respuesta final
	response := services.BuildAPIResponse(users, metadataResponse, paginationResponse)

	c.JSON(http.StatusOK, response)
}

// GetUser obtiene un usuario específico
func GetUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := config.DB.Preload("Role").First(&user, id).Error; err != nil {
		status, response := services.ErrorResponse(services.NewAPIError(
			http.StatusNotFound,
			"NOT_FOUND",
			"Usuario no encontrado",
			"",
			nil,
		))
		c.JSON(status, response)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role.Name,
	})
}

// UpdateUser actualiza un usuario existente
func UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		status, response := services.ErrorResponse(services.NewAPIError(
			http.StatusBadRequest,
			"INVALID_INPUT",
			"Entrada inválida",
			err.Error(),
			nil,
		))
		c.JSON(status, response)
		return
	}

	// Verificar si el usuario tiene permisos para actualizar este usuario
	userId, exists := c.Get("user_id")
	if !exists || (uint(id) != userId.(uint) && c.GetString("user_role") != "admin") {
		status, response := services.ErrorResponse(services.NewAPIError(
			http.StatusForbidden,
			"FORBIDDEN",
			"No tienes permisos para actualizar este usuario",
			"",
			nil,
		))
		c.JSON(status, response)
		return
	}

	var user models.User
	if err := config.DB.Preload("Role").First(&user, id).Error; err != nil {
		status, response := services.ErrorResponse(services.NewAPIError(
			http.StatusNotFound,
			"NOT_FOUND",
			"Usuario no encontrado",
			"",
			nil,
		))
		c.JSON(status, response)
		return
	}

	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		status, response := services.ErrorResponse(services.NewAPIError(
			http.StatusBadRequest,
			"INVALID_INPUT",
			"Entrada inválida",
			err.Error(),
			nil,
		))
		c.JSON(status, response)
		return
	}

	updates := map[string]interface{}{}
	if input.Username != "" {
		updates["username"] = input.Username
	}
	if input.Email != "" {
		updates["email"] = input.Email
	}

	if err := config.DB.Model(&user).Updates(updates).Error; err != nil {
		status, response := services.ErrorResponse(services.NewAPIError(
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"Error interno",
			err.Error(),
			nil,
		))
		c.JSON(status, response)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role.Name,
	})
}

// DeleteUser elimina un usuario
func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	// Solo los administradores pueden eliminar usuarios
	if c.GetString("user_role") != "admin" {
		status, response := services.ErrorResponse(services.NewAPIError(
			http.StatusForbidden,
			"FORBIDDEN",
			"Admin privileges required",
			"",
			nil,
		))
		c.JSON(status, response)
		return
	}

	if err := config.DB.Delete(&models.User{}, id).Error; err != nil {
		status, response := services.ErrorResponse(services.NewAPIError(
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"Error interno",
			err.Error(),
			nil,
		))
		c.JSON(status, response)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
