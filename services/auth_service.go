package services

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	UserId uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type AuthService struct {
	secretKey      string
	expirationHours int
}

func NewAuthService() *AuthService {
	expirationHours, err := strconv.Atoi(os.Getenv("JWT_EXPIRATION_HOURS"))
	if err != nil {
		expirationHours = 24 // valor por defecto
	}

	return &AuthService{
		secretKey:      os.Getenv("JWT_SECRET_KEY"),
		expirationHours: expirationHours,
	}
}

func (s *AuthService) GenerateToken(userId uint, role string) (string, error) {
	claims := TokenClaims{
		UserId: userId,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(s.expirationHours))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("método de firma inesperado")
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("token inválido")
}

func (s *AuthService) RefreshToken(oldTokenString string) (string, error) {
	claims, err := s.ValidateToken(oldTokenString)
	if err != nil {
		return "", err
	}

	// Generar nuevo token
	return s.GenerateToken(claims.UserId, claims.Role)
}

func (s *AuthService) GetUserFromToken(tokenString string) (uint, string, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return 0, "", err
	}

	return claims.UserId, claims.Role, nil
}
