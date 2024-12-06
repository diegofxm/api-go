package services

import (
	"fmt"
	"net/http"
)

// APIError representa un error de la API
type APIError struct {
	Status     int    `json:"-"`              // HTTP status code
	Code       string `json:"code"`           // Error code for client
	Message    string `json:"message"`        // User-friendly error message
	Detail     string `json:"detail"`         // Detailed error message
	Internal   error  `json:"-"`             // Internal error (not exposed)
}

// Error implementa la interfaz error
func (e *APIError) Error() string {
	if e.Internal != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Internal)
	}
	return e.Message
}

// NewAPIError crea un nuevo APIError
func NewAPIError(status int, code, message, detail string, internal error) *APIError {
	return &APIError{
		Status:   status,
		Code:     code,
		Message:  message,
		Detail:   detail,
		Internal: internal,
	}
}

// Common API Errors
var (
	ErrInvalidInput = func(detail string) *APIError {
		return NewAPIError(
			http.StatusBadRequest,
			"INVALID_INPUT",
			"Los datos proporcionados no son válidos",
			detail,
			nil,
		)
	}

	ErrUnauthorized = func(detail string) *APIError {
		return NewAPIError(
			http.StatusUnauthorized,
			"UNAUTHORIZED",
			"No autorizado",
			detail,
			nil,
		)
	}

	ErrForbidden = func(detail string) *APIError {
		return NewAPIError(
			http.StatusForbidden,
			"FORBIDDEN",
			"Acceso denegado",
			detail,
			nil,
		)
	}

	ErrNotFound = func(resource string) *APIError {
		return NewAPIError(
			http.StatusNotFound,
			"NOT_FOUND",
			fmt.Sprintf("%s no encontrado", resource),
			"",
			nil,
		)
	}

	ErrInternal = func(err error) *APIError {
		return NewAPIError(
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"Error interno del servidor",
			"",
			err,
		)
	}

	ErrInvalidCredentials = func() *APIError {
		return NewAPIError(
			http.StatusUnauthorized,
			"INVALID_CREDENTIALS",
			"Credenciales inválidas",
			"",
			nil,
		)
	}
)

// ErrorResponse genera la respuesta de error para la API
func ErrorResponse(err error) (int, interface{}) {
	if apiErr, ok := err.(*APIError); ok {
		response := map[string]interface{}{
			"error": map[string]interface{}{
				"code":    apiErr.Code,
				"message": apiErr.Message,
			},
		}

		if apiErr.Detail != "" {
			response["error"].(map[string]interface{})["detail"] = apiErr.Detail
		}

		return apiErr.Status, response
	}

	// Si no es un APIError, convertirlo en un error interno
	internalErr := ErrInternal(err)
	return internalErr.Status, map[string]interface{}{
		"error": map[string]interface{}{
			"code":    internalErr.Code,
			"message": internalErr.Message,
		},
	}
}
