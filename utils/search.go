package utils

import (
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SearchFilters struct {
	Search string
	Role   string
	SortBy string
	Order  string
}

func ExtractSearchFilters(c *gin.Context) SearchFilters {
	return SearchFilters{
		Search: c.Query("search"),
		Role:   c.Query("role"),
		SortBy: c.Query("sort_by"),
		Order:  c.Query("order"),
	}
}

func ApplySearchFilters(db *gorm.DB, filters SearchFilters) *gorm.DB {
	if filters.Search != "" {
		search := "%" + strings.ToLower(filters.Search) + "%"
		db = db.Where("LOWER(username) LIKE ? OR LOWER(email) LIKE ?", search, search)
	}

	if filters.Role != "" {
		db = db.Where("role = ?", filters.Role)
	}

	if filters.SortBy != "" {
		order := "ASC"
		if strings.ToUpper(filters.Order) == "DESC" {
			order = "DESC"
		}
		db = db.Order(filters.SortBy + " " + order)
	}

	return db
}

// ExtractSortParams extrae los parámetros de ordenamiento de la solicitud
func ExtractSortParams(c *gin.Context) map[string]string {
	sort := make(map[string]string)
	
	// Obtener el parámetro sort
	sortParam := c.Query("sort")
	if sortParam != "" {
		// El formato esperado es "field:direction" (ej: "username:asc")
		parts := strings.Split(sortParam, ":")
		if len(parts) == 2 {
			field := parts[0]
			direction := strings.ToLower(parts[1])
			
			// Validar la dirección
			if direction == "asc" || direction == "desc" {
				sort[field] = direction
			}
		}
	}
	
	return sort
}

// ApplySorting aplica el ordenamiento al query
func ApplySorting(db *gorm.DB, sort map[string]string) *gorm.DB {
	for field, direction := range sort {
		orderStr := field + " " + direction
		db = db.Order(orderStr)
	}
	return db
}
