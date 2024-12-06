package services

import (
	"go-api-orm/models"
	"gorm.io/gorm"
)

// MigrationService maneja las migraciones y datos iniciales
type MigrationService struct {
	db *gorm.DB
}

// NewMigrationService crea una nueva instancia del servicio de migración
func NewMigrationService(db *gorm.DB) *MigrationService {
	return &MigrationService{
		db: db,
	}
}

// SeedDefaultRoles crea los roles por defecto si no existen
func (s *MigrationService) SeedDefaultRoles() error {
	defaultRoles := []models.Role{
		{
			Name:        "admin",
			Description: "Administrador del sistema",
		},
		{
			Name:        "user",
			Description: "Usuario regular",
		},
		{
			Name:        "editor",
			Description: "Editor de contenido",
		},
	}

	for _, role := range defaultRoles {
		var existingRole models.Role
		if err := s.db.Where("name = ?", role.Name).First(&existingRole).Error; err == gorm.ErrRecordNotFound {
			if err := s.db.Create(&role).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
