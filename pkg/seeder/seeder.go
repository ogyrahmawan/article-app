package seeder

import (
	"article-app/internal/domain"

	"gorm.io/gorm"
)

func Seeds(db *gorm.DB) {
	db.Create(&domain.User{
		Email:    "admin@mail.com",
		Password: "Password123",
	})
}
