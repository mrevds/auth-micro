package service

import (
	"auth-micro/internal/auth/entity"
	"context"
)

type UserService interface {
	Register(ctx context.Context, input RegisterInput) (*entity.User, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
}
