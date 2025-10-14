package service

import (
	"auth-micro/internal/auth/entity"
	"context"
)

type UserService interface {
	Register(ctx context.Context, input RegisterInput) (*entity.User, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	Login(ctx context.Context, username, password string) (accessToken string, refreshToken string, err error)
	RefreshAccessToken(ctx context.Context, refreshToken string) (string, error)
	Logout(ctx context.Context, refreshToken string) error
}
