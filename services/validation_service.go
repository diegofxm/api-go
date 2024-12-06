package services

import (
	"regexp"
	"unicode"
)

type ValidationService struct {
	emailRegex    *regexp.Regexp
	usernameRegex *regexp.Regexp
}

func NewValidationService() *ValidationService {
	return &ValidationService{
		emailRegex:    regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`),
		usernameRegex: regexp.MustCompile(`^[a-zA-Z0-9_-]{3,30}$`),
	}
}

// ValidateEmail valida el formato del email
func (s *ValidationService) ValidateEmail(email string) bool {
	return s.emailRegex.MatchString(email)
}

// ValidateUsername valida el formato del nombre de usuario
func (s *ValidationService) ValidateUsername(username string) bool {
	return s.usernameRegex.MatchString(username)
}

// ValidatePassword verifica que la contraseña cumpla con los requisitos mínimos
func (s *ValidationService) ValidatePassword(password string) (bool, string) {
	if len(password) < 8 {
		return false, "La contraseña debe tener al menos 8 caracteres"
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return false, "La contraseña debe contener al menos una mayúscula"
	}
	if !hasLower {
		return false, "La contraseña debe contener al menos una minúscula"
	}
	if !hasNumber {
		return false, "La contraseña debe contener al menos un número"
	}
	if !hasSpecial {
		return false, "La contraseña debe contener al menos un carácter especial"
	}

	return true, ""
}

// ValidateRequired verifica que un campo requerido no esté vacío
func (s *ValidationService) ValidateRequired(field string) bool {
	return field != ""
}

// ValidateLength verifica que un campo tenga una longitud dentro del rango especificado
func (s *ValidationService) ValidateLength(field string, min, max int) bool {
	length := len(field)
	return length >= min && length <= max
}

// ValidateEnum verifica que un valor esté dentro de los valores permitidos
func (s *ValidationService) ValidateEnum(value string, allowedValues []string) bool {
	for _, allowed := range allowedValues {
		if value == allowed {
			return true
		}
	}
	return false
}
