package user

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByUserName(ctx context.Context, username string) (*User, error)
}

type defaultRepository struct {
	database *gorm.DB
}

func NewRepository(database *gorm.DB) Repository {
	return &defaultRepository{
		database: database,
	}
}

func (r *defaultRepository) GetByUserName(ctx context.Context, username string) (*User, error) {
	user := User{}

	timeoutCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := r.database.WithContext(timeoutCtx).Model(&User{}).First(&user, "user_name = ?", username).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *defaultRepository) Create(ctx context.Context, user *User) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := r.database.WithContext(timeoutCtx).Model(&User{}).Create(user).Error; err != nil {
		return err
	}

	return nil
}
