package migrations

import (
	"go-api-orm/models"
	"gorm.io/gorm"
)

// SeedDefaultRoles crea los roles por defecto si no existen
func SeedDefaultRoles(db *gorm.DB) error {
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
		if err := db.Where("name = ?", role.Name).First(&existingRole).Error; err == gorm.ErrRecordNotFound {
			if err := db.Create(&role).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
