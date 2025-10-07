package main

import (
	"auth-micro/client"
	auth "auth-micro/pkg/auth_v1"
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const grpcPort = "50051"

type server struct {
	auth.UnimplementedAuthServer
	db *client.DB
}

func (s *server) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Генерируем ID пользователя
	userID := uuid.New().String()
	now := time.Now()

	// Вставляем пользователя в БД
	query := `
        INSERT INTO users (id, username, name, email, age, bio, password, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `
	name := ""
	if req.Name != nil {
		name = *req.Name
	}
	age := int32(0)
	if req.Age != nil {
		age = *req.Age
	}
	bio := ""
	if req.Bio != nil {
		bio = *req.Bio
	}
	_, err = s.db.Pool.Exec(ctx, query, userID, req.Username, req.Name, req.Email, req.Age, req.Bio, string(hashedPassword), now, now)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	// Возвращаем ответ
	userInfo := &auth.UserInfo{
		Username: req.Username,
		Name:     name,
		Email:    req.Email,
		Age:      age,
		Bio:      bio,
	}

	return &auth.RegisterResponse{
		Id:        userID,
		UserInfo:  userInfo,
		CreatedAt: timestamppb.New(now),
		UpdatedAt: timestamppb.New(now),
	}, nil
}

func main() {
	ctx := context.Background()
	db, err := client.NewDB(ctx)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	auth.RegisterAuthServer(s, &server{db: db})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
