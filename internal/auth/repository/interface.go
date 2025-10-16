package repository

import (
	"auth-micro/internal/auth/entity"
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)

	SaveRefreshToken(ctx context.Context, rt *entity.RefreshToken) error
	GetRefreshToken(ctx context.Context, token string) (*entity.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, token string) error
	RevokeUserRefreshTokens(ctx context.Context, userID string) error
	GetByID(ctx context.Context, id string) (*entity.User, error)
	UpdatePassword(ctx context.Context, userID, hashedPassword string) error
}
