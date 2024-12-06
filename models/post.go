package models

import (
	"strconv"
	"strings"
	"time"
	"unicode"
	"gorm.io/gorm"
)

type Post struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Title     string         `json:"title" gorm:"not null"`
	Slug      string         `json:"slug" gorm:"type:varchar(255);uniqueIndex;not null"`
	Content   string         `json:"content"`
	AuthorID  uint           `json:"author_id" gorm:"not null"`
	Author    User           `json:"author" gorm:"foreignKey:AuthorID"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// generateSlug generates a URL-friendly slug from a title
func generateSlug(title string) string {
	// Convert to lowercase
	slug := strings.ToLower(title)
	
	// Replace special characters with their closest ASCII equivalent
	specialChars := map[rune]string{
		'á': "a", 'é': "e", 'í': "i", 'ó': "o", 'ú': "u",
		'ñ': "n", 'ü': "u",
		'&': "and",
	}
	
	var result strings.Builder
	for _, r := range slug {
		if replacement, ok := specialChars[r]; ok {
			result.WriteString(replacement)
		} else if unicode.IsLetter(r) || unicode.IsNumber(r) || r == ' ' {
			result.WriteRune(r)
		}
	}
	
	// Replace spaces with hyphens
	slug = strings.ReplaceAll(result.String(), " ", "-")
	
	// Remove any remaining non-alphanumeric characters (except hyphens)
	slug = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, slug)
	
	// Remove multiple consecutive hyphens
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}
	
	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")
	
	return slug
}

// BeforeCreate is a GORM hook that runs before creating a record
func (p *Post) BeforeCreate(tx *gorm.DB) error {
	if p.Slug == "" {
		baseSlug := generateSlug(p.Title)
		slug := baseSlug
		counter := 1

		// Check if slug exists and generate a unique one
		for {
			var count int64
			tx.Model(&Post{}).Where("slug = ?", slug).Count(&count)
			if count == 0 {
				break
			}
			slug = baseSlug + "-" + strconv.Itoa(counter)
			counter++
		}
		p.Slug = slug
	}
	return nil
}

// BeforeUpdate is a GORM hook that runs before updating a record
func (p *Post) BeforeUpdate(tx *gorm.DB) error {
	// If title changed and slug wasn't manually specified, update it
	var oldPost Post
	if err := tx.First(&oldPost, p.ID).Error; err != nil {
		return err
	}

	if oldPost.Title != p.Title && p.Slug == oldPost.Slug {
		baseSlug := generateSlug(p.Title)
		slug := baseSlug
		counter := 1

		// Check if slug exists and generate a unique one
		for {
			var count int64
			tx.Model(&Post{}).Where("slug = ? AND id != ?", slug, p.ID).Count(&count)
			if count == 0 {
				break
			}
			slug = baseSlug + "-" + strconv.Itoa(counter)
			counter++
		}
		p.Slug = slug
	}
	return nil
}
