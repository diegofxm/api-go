package services

// SearchField representa un campo de búsqueda permitido
type SearchField struct {
	Name        string   `json:"Name"`
	Type        string   `json:"Type"`
	Description string   `json:"Description"`
	Operators   []string `json:"Operators"`
}

// SortField representa un campo de ordenamiento permitido
type SortField struct {
	Name        string `json:"Name"`
	Description string `json:"Description"`
}

// MetadataResponse representa la estructura de metadatos
type MetadataResponse struct {
	AllowedSearch  []SearchField         `json:"allowed_search"`
	AllowedSort    []SortField          `json:"allowed_sort"`
	AppliedSearch  []map[string]string   `json:"applied_search"`
	AppliedSort    map[string]string     `json:"applied_sort"`
}

// BuildMetadataResponse construye la respuesta de metadatos
func BuildMetadataResponse(searchFields []SearchField, sortFields []SortField, appliedSearch []map[string]string, appliedSort map[string]string) MetadataResponse {
	return MetadataResponse{
		AllowedSearch:  searchFields,
		AllowedSort:    sortFields,
		AppliedSearch:  appliedSearch,
		AppliedSort:    appliedSort,
	}
}

// GetDefaultUserSearchFields retorna los campos de búsqueda predefinidos para usuarios
func GetDefaultUserSearchFields() []SearchField {
	return []SearchField{
		{
			Name:        "username",
			Type:        "string",
			Description: "Username del usuario",
			Operators:   []string{"eq", "like", "nlike"},
		},
		{
			Name:        "email",
			Type:        "string",
			Description: "Email del usuario",
			Operators:   []string{"eq", "like", "nlike"},
		},
		{
			Name:        "role",
			Type:        "string",
			Description: "Rol del usuario (admin o user)",
			Operators:   []string{"eq", "ne", "in", "nin"},
		},
		{
			Name:        "created_at",
			Type:        "date",
			Description: "Fecha de creación del usuario",
			Operators:   []string{"gt", "gte", "lt", "lte"},
		},
	}
}

// GetDefaultUserSortFields retorna los campos de ordenamiento predefinidos para usuarios
func GetDefaultUserSortFields() []SortField {
	return []SortField{
		{
			Name:        "username",
			Description: "Ordenar por username",
		},
		{
			Name:        "email",
			Description: "Ordenar por email",
		},
		{
			Name:        "role",
			Description: "Ordenar por rol",
		},
		{
			Name:        "created_at",
			Description: "Ordenar por fecha de creación",
		},
	}
}
