package handler

import (
	"auth-micro/internal/auth/service"
	auth "auth-micro/pkg/auth_v1"
	"context"

	"google.golang.org/grpc/codes"
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
