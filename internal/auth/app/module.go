package app

import (
	"go.uber.org/fx"

	handlerAuth "auth-micro/internal/auth/handler"
	repoPostgres "auth-micro/internal/auth/repository/postgres"
	serviceAuth "auth-micro/internal/auth/service"
)

var Module = fx.Module("app",
	// Repository layer
	fx.Provide(repoPostgres.NewUserRepo),

	// Service layer
	fx.Provide(serviceAuth.NewUserService),

	// Handler layer
	fx.Provide(handlerAuth.NewGRPCHandler),
)
