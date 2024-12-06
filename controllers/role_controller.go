package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-api-orm/config"
	"go-api-orm/models"
	"go-api-orm/services"
)

// CreateRoleInput representa los datos necesarios para crear un rol
type CreateRoleInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// UpdateRoleInput representa los datos que se pueden actualizar de un rol
type UpdateRoleInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CreateRole crea un nuevo rol
func CreateRole(c *gin.Context) {
	var input CreateRoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		status, response := services.ErrorResponse(services.ErrInvalidInput(err.Error()))
		c.JSON(status, response)
		return
	}

	role := models.Role{
		Name:        input.Name,
		Description: input.Description,
	}

	if err := config.DB.Create(&role).Error; err != nil {
		status, response := services.ErrorResponse(services.ErrInternal(err))
		c.JSON(status, response)
		return
	}

	c.JSON(http.StatusCreated, role)
}

// GetRoles obtiene todos los roles con paginación y filtros
func GetRoles(c *gin.Context) {
	var roles []models.Role
	
	// Obtener parámetros de paginación y búsqueda
	pagination := services.GeneratePaginationFromRequest(c)
	searchFilters := services.ExtractSearchParams(c)
	sortParams := services.ExtractSortParams(c)
	
	// Aplicar filtros y paginación
	db := config.DB
	db = services.ApplySearchFilters(db, searchFilters)
	db = services.ApplySorting(db, sortParams)
	
	err := db.Scopes(services.Paginate(roles, &pagination, db)).Find(&roles).Error
	if err != nil {
		status, response := services.ErrorResponse(services.ErrInternal(err))
		c.JSON(status, response)
		return
	}

	// Construir los componentes de la respuesta
	paginationResponse := services.BuildPaginationResponse(c, pagination.Page, pagination.Limit, pagination.TotalRows, pagination.TotalPages)
	
	metadataResponse := services.BuildMetadataResponse(
		[]services.SearchField{
			{
				Name:        "name",
				Type:        "string",
				Description: "Nombre del rol",
				Operators:   []string{"eq", "like", "nlike"},
			},
			{
				Name:        "description",
				Type:        "string",
				Description: "Descripción del rol",
				Operators:   []string{"like", "nlike"},
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
				Name:        "name",
				Description: "Ordenar por nombre",
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
	response := services.BuildAPIResponse(roles, metadataResponse, paginationResponse)

	c.JSON(http.StatusOK, response)
}

// GetRole obtiene un rol por su ID
func GetRole(c *gin.Context) {
	id := c.Param("id")
	
	var role models.Role
	if err := config.DB.First(&role, id).Error; err != nil {
		status, response := services.ErrorResponse(services.ErrNotFound("Role"))
		c.JSON(status, response)
		return
	}

	c.JSON(http.StatusOK, role)
}

// UpdateRole actualiza un rol existente
func UpdateRole(c *gin.Context) {
	id := c.Param("id")
	
	var role models.Role
	if err := config.DB.First(&role, id).Error; err != nil {
		status, response := services.ErrorResponse(services.ErrNotFound("Role"))
		c.JSON(status, response)
		return
	}

	var input UpdateRoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		status, response := services.ErrorResponse(services.ErrInvalidInput(err.Error()))
		c.JSON(status, response)
		return
	}

	// Actualizar solo los campos proporcionados
	updates := map[string]interface{}{}
	if input.Name != "" {
		updates["name"] = input.Name
	}
	if input.Description != "" {
		updates["description"] = input.Description
	}

	if err := config.DB.Model(&role).Updates(updates).Error; err != nil {
		status, response := services.ErrorResponse(services.ErrInternal(err))
		c.JSON(status, response)
		return
	}

	c.JSON(http.StatusOK, role)
}

// DeleteRole elimina un rol
func DeleteRole(c *gin.Context) {
	id := c.Param("id")
	
	var role models.Role
	if err := config.DB.First(&role, id).Error; err != nil {
		status, response := services.ErrorResponse(services.ErrNotFound("Role"))
		c.JSON(status, response)
		return
	}

	if err := config.DB.Delete(&role).Error; err != nil {
		status, response := services.ErrorResponse(services.ErrInternal(err))
		c.JSON(status, response)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role deleted successfully"})
}
