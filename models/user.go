package models

import (
	"time"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"uniqueIndex:idx_username,length:255;not null;size:255"`
	Email     string         `json:"email" gorm:"uniqueIndex:idx_email,length:255;not null;size:255"`
	Password  string         `json:"-" gorm:"not null"`
	RoleID    uint           `json:"role_id" gorm:"not null"`
	Role      Role           `json:"role" gorm:"foreignKey:RoleID"`
	Posts     []Post         `json:"posts,omitempty" gorm:"foreignKey:AuthorID"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// BeforeCreate es un hook que se ejecuta antes de crear un usuario
func (u *User) BeforeCreate(tx *gorm.DB) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword verifica si la contraseña proporcionada coincide con la almacenada
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

// TableName especifica el nombre de la tabla para GORM
func (User) TableName() string {
	return "users"
}
