package user

import (
	"context"
	"errors"

	"github.com/nkyizbay/shop-project/internal/auth"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(ctx context.Context, user *User) error
	Login(ctx context.Context, credentials auth.Credentials) (*User, error)
}

type defaultService struct {
	userRepo Repository
}

func NewUserService(repository Repository) Service {
	return &defaultService{userRepo: repository}
}

func (s *defaultService) Register(ctx context.Context, user *User) error {
	err := s.userRepo.Create(ctx, user)

	if err != nil {
		return err
	}

	return nil
}

func (s *defaultService) Login(ctx context.Context, credentials auth.Credentials) (*User, error) {
	user, err := s.userRepo.GetByUserName(ctx, credentials.UserName)
	if err != nil {
		return nil, err
	}

	if s.isNotEqualHashAndPassword(user.Password, credentials.Password) {
		return nil, errors.New("invalid password")
	}

	return user, nil
}

func (s *defaultService) isNotEqualHashAndPassword(hashPassword string, password string) bool {
	return !s.isEqualHashAndPassword(hashPassword, password)
}

func (s *defaultService) isEqualHashAndPassword(hashPassword string, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password)); err != nil {
		return false
	}
	return true
}
