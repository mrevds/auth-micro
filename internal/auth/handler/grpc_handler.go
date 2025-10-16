package handler

import (
	"auth-micro/internal/auth/service"
	auth "auth-micro/pkg/auth_v1"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type grpcHandler struct {
	auth.UnimplementedAuthServer
	userService service.UserService
}

func NewGRPCHandler(s service.UserService) auth.AuthServer {
	return &grpcHandler{userService: s}
}

func (h *grpcHandler) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	user, err := h.userService.Register(ctx, service.RegisterInput{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
		Age:      req.Age,
		Bio:      req.Bio,
	})
	if err != nil {
		return nil, err
	}

	return &auth.RegisterResponse{
		Id: user.ID,
		UserInfo: &auth.UserInfo{
			Id:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Name:      user.Name,
			Age:       user.Age,
			Bio:       user.Bio,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		},
	}, nil
}

func (h *grpcHandler) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	if req.Username == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "username and password are required")
	}

	accessToken, refreshToken, err := h.userService.Login(ctx, req.Username, req.Password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	return &auth.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (h *grpcHandler) Logout(ctx context.Context, req *auth.LogoutRequest) (*auth.LogoutResponse, error) {
	if req.RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh token is required")
	}

	if err := h.userService.Logout(ctx, req.RefreshToken); err != nil {
		return nil, status.Error(codes.Internal, "failed to logout")
	}

	return &auth.LogoutResponse{
		Success: true,
		Message: "Logged out successfully",
	}, nil
}

func (h *grpcHandler) ChangePassword(ctx context.Context, req *auth.ChangePasswordRequest) (*auth.ChangePasswordResponse, error) {
	// Валидация
	if req.NewPassword == "" || req.OldPassword == "" {
		return nil, status.Error(codes.InvalidArgument, "current and new passwords are required")
	}

	// Получить токен из metadata (Handler НЕ знает что внутри токена!)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing token")
	}

	tokens := md.Get("authorization")
	if len(tokens) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing token")
	}

	// Просто передаем токен в Service (Handler не парсит его!)
	err := h.userService.ChangePassword(ctx, tokens[0], req.OldPassword, req.NewPassword)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &auth.ChangePasswordResponse{
		Success: true,
		Message: "Password changed successfully",
	}, nil
}
