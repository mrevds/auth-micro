package auth

import (
	"auth-micro/client"
	"auth-micro/internal/auth/handler"
	"auth-micro/internal/auth/repository"
	"auth-micro/internal/auth/repository/postgres"
	"auth-micro/internal/auth/service"
)

type Module struct {
	Handler handler.AuthHandler // gRPC интерфейс
	Service service.UserService
	Repo    repository.UserRepository
}

func NewModule(db *client.DB) *Module {
	repo := postgres.NewUserRepo(db)
	svc := service.NewUserService(repo)
	h := handler.NewGRPCHandler(svc)

	return &Module{
		Handler: h,
		Service: svc,
		Repo:    repo,
	}
}
