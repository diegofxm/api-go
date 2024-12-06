package services

import (
	"html"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type TransformService struct {
	dateFormat    string
	slugSeparator string
	slugRegex     *regexp.Regexp
}

func NewTransformService() *TransformService {
	return &TransformService{
		dateFormat:    "2006-01-02T15:04:05.0000000-07:00",
		slugSeparator: "-",
		slugRegex:     regexp.MustCompile("[^a-z0-9]+"),
	}
}

// FormatDateTime formatea una fecha al formato estándar de la API
func (s *TransformService) FormatDateTime(t time.Time) string {
	return t.Format(s.dateFormat)
}

// ParseDateTime parsea una fecha en formato string al tipo time.Time
func (s *TransformService) ParseDateTime(dateStr string) (time.Time, error) {
	return time.Parse(s.dateFormat, dateStr)
}

// SanitizeString limpia una cadena de texto de posibles caracteres maliciosos
func (s *TransformService) SanitizeString(input string) string {
	// Escapar HTML
	escaped := html.EscapeString(input)
	// Remover caracteres de control
	return strings.Map(func(r rune) rune {
		if r < 32 || r == 127 {
			return -1
		}
		return r
	}, escaped)
}

// TrimSpaces elimina espacios en blanco al inicio y final
func (s *TransformService) TrimSpaces(input string) string {
	return strings.TrimSpace(input)
}

// NormalizeEmail normaliza un email (convierte a minúsculas y elimina espacios)
func (s *TransformService) NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// ToSnakeCase convierte una cadena de CamelCase a snake_case
func (s *TransformService) ToSnakeCase(str string) string {
	var result strings.Builder
	for i, r := range str {
		if i > 0 && unicode.IsUpper(r) {
			result.WriteRune('_')
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}

// ToCamelCase convierte una cadena de snake_case a CamelCase
func (s *TransformService) ToCamelCase(str string) string {
	words := strings.Split(str, "_")
	for i := 1; i < len(words); i++ {
		words[i] = strings.Title(words[i])
	}
	return strings.Join(words, "")
}

// Truncate acorta una cadena a una longitud máxima
func (s *TransformService) Truncate(str string, maxLength int) string {
	if len(str) <= maxLength {
		return str
	}
	return str[:maxLength-3] + "..."
}

// GenerateSlug genera un slug a partir de un texto
// Ejemplo: "Hello World!" -> "hello-world"
func (s *TransformService) GenerateSlug(text string) string {
	// Convertir a minúsculas
	text = strings.ToLower(text)

	// Reemplazar caracteres especiales y espacios con el separador
	text = s.slugRegex.ReplaceAllString(text, s.slugSeparator)

	// Eliminar separadores del inicio y final
	text = strings.Trim(text, s.slugSeparator)

	return text
}

// GenerateUniqueSlug genera un slug único agregando un sufijo si es necesario
func (s *TransformService) GenerateUniqueSlug(text string, existingSlug func(string) bool) string {
	baseSlug := s.GenerateSlug(text)
	slug := baseSlug
	counter := 1

	// Mientras el slug exista, agregar un sufijo numérico
	for existingSlug(slug) {
		slug = baseSlug + s.slugSeparator + strconv.Itoa(counter)
		counter++
	}

	return slug
}

// SetSlugSeparator cambia el separador usado en los slugs
func (s *TransformService) SetSlugSeparator(separator string) {
	s.slugSeparator = separator
}

// ValidateSlug verifica si un slug es válido
func (s *TransformService) ValidateSlug(slug string) bool {
	// Un slug válido solo debe contener letras minúsculas, números y el separador
	validSlugRegex := regexp.MustCompile("^[a-z0-9]+(?:" + regexp.QuoteMeta(s.slugSeparator) + "[a-z0-9]+)*$")
	return validSlugRegex.MatchString(slug)
}

// NormalizeSlug normaliza un slug existente
func (s *TransformService) NormalizeSlug(slug string) string {
	return s.GenerateSlug(slug)
}
