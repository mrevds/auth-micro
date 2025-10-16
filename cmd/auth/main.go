package main

import (
    "context"
    "fmt"
    "log"
    "net"

    "go.uber.org/fx"
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"

    "auth-micro/client"
    "auth-micro/internal/auth/app"
    "auth-micro/internal/auth/config"
    pb "auth-micro/pkg/auth_v1"
)

func main() {
    fx.New(
        // Провайдеры
        fx.Provide(
            config.Load,           // Загрузка конфигурации
            client.NewDB,          // БД с lifecycle
            newGRPCServer,         // gRPC сервер
        ),

        // Модули
        app.Module,

        // Lifecycle для gRPC сервера
        fx.Invoke(registerGRPCServer),
    ).Run()
}

func newGRPCServer() *grpc.Server {
    return grpc.NewServer()
}

func registerGRPCServer(
    lc fx.Lifecycle,
    grpcServer *grpc.Server,
    handler pb.AuthServer,
    cfg *config.Config,
) {
    lc.Append(fx.Hook{
        OnStart: func(ctx context.Context) error {
            lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Server.GRPCPort))
            if err != nil {
                return fmt.Errorf("failed to listen: %w", err)
            }

            pb.RegisterAuthServer(grpcServer, handler)
            reflection.Register(grpcServer)

            go func() {
                log.Printf("✅ gRPC server listening on port %s", cfg.Server.GRPCPort)
                if err := grpcServer.Serve(lis); err != nil {
                    log.Fatalf("Failed to serve: %v", err)
                }
            }()

            return nil
        },
        OnStop: func(ctx context.Context) error {
            log.Println("🛑 Stopping gRPC server...")
            grpcServer.GracefulStop()
            log.Println("✅ gRPC server stopped")
            return nil
        },
    })
}