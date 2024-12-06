package services

import (
	"os"
	"strings"
)

// APIResponse representa la estructura de respuesta general de la API
type APIResponse struct {
	Data       interface{}       `json:"data"`
	Metadata   *MetadataResponse  `json:"metadata,omitempty"`
	Pagination *PaginationResponse `json:"pagination,omitempty"`
}

// isTrue verifica si una variable de entorno es "true"
func isTrue(envVar string) bool {
	return strings.ToLower(os.Getenv(envVar)) == "true"
}

// BuildAPIResponse construye la respuesta completa de la API
func BuildAPIResponse(data interface{}, metadata MetadataResponse, pagination PaginationResponse) APIResponse {
	response := APIResponse{
		Data: data,
	}

	// Agregar metadata si está habilitada
	if isTrue("SHOW_METADATA") {
		response.Metadata = &metadata
	}

	// Agregar paginación si está habilitada
	if isTrue("SHOW_PAGINATION") {
		response.Pagination = &pagination
	}

	return response
}
