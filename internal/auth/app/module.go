package app

import (
    "go.uber.org/fx"

    "auth-micro/internal/auth/handler"
    repoPostgres "auth-micro/internal/auth/repository/postgres"
    serviceAuth "auth-micro/internal/auth/service"
    "auth-micro/internal/auth/utils"
)

var Module = fx.Module("app",
    fx.Provide(repoPostgres.NewUserRepo),
    fx.Provide(utils.NewJWTManager),
    fx.Provide(serviceAuth.NewUserService),
    fx.Provide(handler.NewGRPCHandler),
)