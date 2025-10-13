package service

import (
	"auth-micro/internal/auth/entity"
	"auth-micro/internal/auth/repository"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

type RegisterInput struct {
	Username string
	Email    string
	Password string
	Name     *string
	Age      *int32
	Bio      *string
}

func (s *userService) Register(ctx context.Context, input RegisterInput) (*entity.User, error) {
	existing, _ := s.repo.GetByUsername(ctx, input.Username)
	if existing != nil {
		return nil, fmt.Errorf("username already taken")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		ID:        uuid.NewString(),
		Username:  input.Username,
		Email:     input.Email,
		Name:      getString(input.Name),
		Age:       getInt32(input.Age),
		Bio:       getString(input.Bio),
		Password:  string(hashed),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	return s.repo.GetByUsername(ctx, username)
}

func getString(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}
func getInt32(v *int32) int32 {
	if v == nil {
		return 0
	}
	return *v
}
