package services

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Pagination struct {
	Limit      int    `json:"limit,omitempty;query:limit"`
	Page       int    `json:"page,omitempty;query:page"`
	TotalRows  int64  `json:"total_rows"`
	TotalPages int    `json:"total_pages"`
}

type PaginationLinks struct {
	First string `json:"first,omitempty"`
	Last  string `json:"last,omitempty"`
	Prev  string `json:"prev,omitempty"`
	Next  string `json:"next,omitempty"`
}

type PaginationResponse struct {
	CurrentPage int             `json:"current_page"`
	PerPage     int            `json:"per_page"`
	TotalItems  int64          `json:"total_items"`
	TotalPages  int            `json:"total_pages"`
	Links       PaginationLinks `json:"links"`
}

// GeneratePaginationFromRequest genera la paginación a partir de los parámetros de la solicitud
func GeneratePaginationFromRequest(c *gin.Context) Pagination {
	// Limit -> perPage
	limit := 10
	page := 1
	query := c.Request.URL.Query()
	for key, value := range query {
		queryValue := value[len(value)-1]
		switch key {
		case "limit":
			limit, _ = strconv.Atoi(queryValue)
		case "page":
			page, _ = strconv.Atoi(queryValue)
		}
	}
	return Pagination{
		Limit: limit,
		Page:  page,
	}
}

// Paginate aplica la paginación al query
func Paginate(value interface{}, pagination *Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	db.Model(value).Count(&totalRows)

	pagination.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))
	pagination.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset((pagination.Page - 1) * pagination.Limit).Limit(pagination.Limit)
	}
}

// BuildPaginationLinks construye los enlaces de paginación
func BuildPaginationLinks(baseURL string, currentPage, totalPages int) PaginationLinks {
	links := PaginationLinks{
		First: baseURL + "?page=1",
		Last:  baseURL + "?page=" + strconv.Itoa(totalPages),
	}

	if currentPage > 1 {
		links.Prev = baseURL + "?page=" + strconv.Itoa(currentPage-1)
	}
	if currentPage < totalPages {
		links.Next = baseURL + "?page=" + strconv.Itoa(currentPage+1)
	}

	return links
}

// BuildPaginationResponse construye la respuesta de paginación completa
func BuildPaginationResponse(c *gin.Context, page, limit int, totalRows int64, totalPages int) PaginationResponse {
	baseURL := c.Request.URL.Path
	links := BuildPaginationLinks(baseURL, page, totalPages)

	return PaginationResponse{
		CurrentPage: page,
		PerPage:     limit,
		TotalItems:  totalRows,
		TotalPages:  totalPages,
		Links:       links,
	}
}
