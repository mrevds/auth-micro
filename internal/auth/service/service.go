package service

import (
	"auth-micro/internal/auth/entity"
	"auth-micro/internal/auth/repository"
	"auth-micro/internal/auth/utils"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo       repository.UserRepository
	jwtManager *utils.JWTManager
}

func NewUserService(repo repository.UserRepository, jwtManager *utils.JWTManager) UserService {
	return &userService{
		repo:       repo,
		jwtManager: jwtManager,
	}
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

func (s *userService) Login(ctx context.Context, username, password string) (string, string, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return "", "", fmt.Errorf("database error: %w", err)
	}
	if user == nil {
		return "", "", fmt.Errorf("invalid credentials")
	}

	if err := utils.CheckPasswordHash(password, user.Password); err != nil {
		return "", "", fmt.Errorf("invalid credentials")
	}

	// Использование JWTManager вместо прямых вызовов utils
	accessToken, err := s.jwtManager.GenerateToken(user.ID)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	rt := &entity.RefreshToken{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		CreatedAt: time.Now(),
		Revoked:   false,
	}
	if err := s.repo.SaveRefreshToken(ctx, rt); err != nil {
		return "", "", fmt.Errorf("failed to save refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

func (s *userService) RefreshAccessToken(ctx context.Context, refreshToken string) (string, error) {
	// Использование JWTManager
	claims, err := s.jwtManager.ValidateToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token")
	}

	if claims.Type != "refresh" {
		return "", fmt.Errorf("invalid token type")
	}

	rt, err := s.repo.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", fmt.Errorf("database error: %w", err)
	}
	if rt == nil || rt.Revoked {
		return "", fmt.Errorf("refresh token revoked or not found")
	}

	if time.Now().After(rt.ExpiresAt) {
		return "", fmt.Errorf("refresh token expired")
	}

	newAccessToken, err := s.jwtManager.GenerateToken(claims.UserID)
	if err != nil {
		return "", fmt.Errorf("failed to generate new access token: %w", err)
	}

	return newAccessToken, nil
}

func (s *userService) ChangePassword(ctx context.Context, token, oldPassword, newPassword string) error {
	// Использование JWTManager
	claims, err := s.jwtManager.ValidateToken(token)
	if err != nil {
		return fmt.Errorf("invalid token: %w", err)
	}

	user, err := s.repo.GetByID(ctx, claims.UserID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	if err := utils.CheckPasswordHash(oldPassword, user.Password); err != nil {
		return fmt.Errorf("current password is incorrect")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	return s.repo.UpdatePassword(ctx, claims.UserID, hashedPassword)
}

func (s *userService) Logout(ctx context.Context, refreshToken string) error {
	return s.repo.RevokeRefreshToken(ctx, refreshToken)
}

func (s *userService) GetUserInfo(ctx context.Context, username, token string) (*entity.User, error) {
	claims, err := s.jwtManager.ValidateToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if claims.Type != "access" {
		return nil, fmt.Errorf("invalid token type")
	}

	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Проверяем, что запрашиваемый пользователь совпадает с пользователем из токена
	if user.ID != claims.UserID {
		return nil, fmt.Errorf("unauthorized access")
	}

	return user, nil
}

func (s *userService) GetUserByID(ctx context.Context, userID string) (*entity.User, error) {
	return s.repo.GetByID(ctx, userID)
}
