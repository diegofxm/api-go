package services

import (
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SearchParams representa los parámetros de búsqueda
type SearchParams struct {
	Field    string
	Operator string
	Value    string
}

// SortParams representa los parámetros de ordenamiento
type SortParams struct {
	Field     string
	Direction string
}

// ExtractSearchParams extrae los parámetros de búsqueda de la solicitud
func ExtractSearchParams(c *gin.Context) []map[string]string {
	var filters []map[string]string
	search := c.QueryArray("search")

	for _, s := range search {
		parts := strings.Split(s, ":")
		if len(parts) == 3 {
			filter := map[string]string{
				"field":    parts[0],
				"operator": parts[1],
				"value":    parts[2],
			}
			filters = append(filters, filter)
		}
	}

	return filters
}

// ExtractSortParams extrae los parámetros de ordenamiento de la solicitud
func ExtractSortParams(c *gin.Context) map[string]string {
	sort := make(map[string]string)
	
	sortParam := c.Query("sort")
	if sortParam != "" {
		parts := strings.Split(sortParam, ":")
		if len(parts) == 2 {
			field := parts[0]
			direction := strings.ToLower(parts[1])
			
			if direction == "asc" || direction == "desc" {
				sort[field] = direction
			}
		}
	}
	
	return sort
}

// ApplySearchFilters aplica los filtros de búsqueda al query
func ApplySearchFilters(db *gorm.DB, filters []map[string]string) *gorm.DB {
	for _, filter := range filters {
		field := filter["field"]
		operator := filter["operator"]
		value := filter["value"]

		switch operator {
		case "eq":
			db = db.Where(field+" = ?", value)
		case "ne":
			db = db.Where(field+" != ?", value)
		case "like":
			db = db.Where(field+" LIKE ?", "%"+value+"%")
		case "nlike":
			db = db.Where(field+" NOT LIKE ?", "%"+value+"%")
		case "in":
			values := strings.Split(value, ",")
			db = db.Where(field+" IN ?", values)
		case "nin":
			values := strings.Split(value, ",")
			db = db.Where(field+" NOT IN ?", values)
		case "gt":
			db = db.Where(field+" > ?", value)
		case "gte":
			db = db.Where(field+" >= ?", value)
		case "lt":
			db = db.Where(field+" < ?", value)
		case "lte":
			db = db.Where(field+" <= ?", value)
		}
	}

	return db
}

// ApplySorting aplica el ordenamiento al query
func ApplySorting(db *gorm.DB, sort map[string]string) *gorm.DB {
	for field, direction := range sort {
		orderStr := field + " " + direction
		db = db.Order(orderStr)
	}
	return db
}
