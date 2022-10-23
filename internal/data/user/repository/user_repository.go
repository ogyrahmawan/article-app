package repository

import (
	"article-app/internal/domain"
	"context"

	"gorm.io/gorm"
)

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{
		DB: db,
	}
}

func (ur userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var entity domain.User
	err := ur.DB.WithContext(ctx).First(&entity, "email =?", email).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}
