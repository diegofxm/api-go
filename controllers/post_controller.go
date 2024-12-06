package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-api-orm/config"
	"go-api-orm/models"
	"go-api-orm/services"
)

type CreatePostInput struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	Slug    string `json:"slug"` // opcional
}

type UpdatePostInput struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Slug    string `json:"slug"`
}

// CreatePost crea un nuevo post
func CreatePost(c *gin.Context) {
	var input CreatePostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		status, response := services.ErrorResponse(services.ErrInvalidInput(err.Error()))
		c.JSON(status, response)
		return
	}

	// Obtener el ID del usuario del token
	userId, ok := c.Get("userId")
	if !ok {
		status, response := services.ErrorResponse(services.ErrInvalidInput("No se encontró el ID del usuario"))
		c.JSON(status, response)
		return
	}

	// Crear el post
	post := models.Post{
		Title:    input.Title,
		Content:  input.Content,
		Slug:     input.Slug, // Si está vacío, el hook BeforeCreate generará uno
		AuthorID: userId.(uint), // Obtener el ID del usuario del token
	}

	if err := config.DB.Create(&post).Error; err != nil {
		status, response := services.ErrorResponse(services.ErrInternal(err))
		c.JSON(status, response)
		return
	}

	c.JSON(http.StatusCreated, post)
}

// GetPostBySlug obtiene un post por su slug
func GetPostBySlug(c *gin.Context) {
	slug := c.Param("slug")
	
	var post models.Post
	if err := config.DB.Where("slug = ?", slug).First(&post).Error; err != nil {
		status, response := services.ErrorResponse(services.ErrNotFound("Post"))
		c.JSON(status, response)
		return
	}

	c.JSON(http.StatusOK, post)
}

// GetPosts obtiene todos los posts con paginación y filtros
func GetPosts(c *gin.Context) {
	var posts []models.Post
	
	// Obtener parámetros de paginación y búsqueda
	pagination := services.GeneratePaginationFromRequest(c)
	searchFilters := services.ExtractSearchParams(c)
	sortParams := services.ExtractSortParams(c)
	
	// Aplicar filtros y paginación
	db := config.DB
	db = services.ApplySearchFilters(db, searchFilters)
	db = services.ApplySorting(db, sortParams)
	
	err := db.Scopes(services.Paginate(posts, &pagination, db)).Find(&posts).Error
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
				Name:        "title",
				Type:        "string",
				Description: "Título del post",
				Operators:   []string{"eq", "like", "nlike"},
			},
			{
				Name:        "slug",
				Type:        "string",
				Description: "Slug del post",
				Operators:   []string{"eq"},
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
				Name:        "title",
				Description: "Ordenar por título",
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
	response := services.BuildAPIResponse(posts, metadataResponse, paginationResponse)

	c.JSON(http.StatusOK, response)
}

// UpdatePost actualiza un post existente
func UpdatePost(c *gin.Context) {
	slug := c.Param("slug")
	
	var post models.Post
	if err := config.DB.Where("slug = ?", slug).First(&post).Error; err != nil {
		status, response := services.ErrorResponse(services.ErrNotFound("Post"))
		c.JSON(status, response)
		return
	}

	var input UpdatePostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		status, response := services.ErrorResponse(services.ErrInvalidInput(err.Error()))
		c.JSON(status, response)
		return
	}

	// Actualizar solo los campos proporcionados
	updates := map[string]interface{}{}
	if input.Title != "" {
		updates["title"] = input.Title
	}
	if input.Content != "" {
		updates["content"] = input.Content
	}
	if input.Slug != "" {
		// Validar el slug personalizado
		transformService := services.NewTransformService()
		if !transformService.ValidateSlug(input.Slug) {
			status, response := services.ErrorResponse(services.ErrInvalidInput("Slug inválido"))
			c.JSON(status, response)
			return
		}
		// Verificar que el slug no exista
		var count int64
		if err := config.DB.Model(&models.Post{}).Where("slug = ? AND id != ?", input.Slug, post.ID).Count(&count).Error; err != nil {
			status, response := services.ErrorResponse(services.ErrInternal(err))
			c.JSON(status, response)
			return
		}
		if count > 0 {
			status, response := services.ErrorResponse(services.ErrInvalidInput("Slug ya existe"))
			c.JSON(status, response)
			return
		}
		updates["slug"] = input.Slug
	}

	if err := config.DB.Model(&post).Updates(updates).Error; err != nil {
		status, response := services.ErrorResponse(services.ErrInternal(err))
		c.JSON(status, response)
		return
	}

	c.JSON(http.StatusOK, post)
}

// DeletePost elimina un post
func DeletePost(c *gin.Context) {
	slug := c.Param("slug")
	
	var post models.Post
	if err := config.DB.Where("slug = ?", slug).First(&post).Error; err != nil {
		status, response := services.ErrorResponse(services.ErrNotFound("Post"))
		c.JSON(status, response)
		return
	}

	if err := config.DB.Delete(&post).Error; err != nil {
		status, response := services.ErrorResponse(services.ErrInternal(err))
		c.JSON(status, response)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post eliminado correctamente"})
}
