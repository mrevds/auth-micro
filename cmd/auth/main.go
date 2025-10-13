package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/joho/godotenv"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"auth-micro/client"
	"auth-micro/internal/auth/app"
	pb "auth-micro/pkg/auth_v1"
)

const grpcPort = "50051"

func main() {
	_ = godotenv.Load()

	fx.New(
		fx.Provide(
			client.NewDB,  // БД с lifecycle
			newGRPCServer, // gRPC сервер
		),
		app.Module,
		fx.Invoke(startServer),
	).Run()
}

// newGRPCServer создает gRPC сервер с зарегистрированным auth handler
func newGRPCServer(authHandler pb.AuthServer) *grpc.Server {
	srv := grpc.NewServer()
	pb.RegisterAuthServer(srv, authHandler)
	reflection.Register(srv)
	return srv
}

// startServer запускает gRPC сервер
func startServer(lifecycle fx.Lifecycle, srv *grpc.Server) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
			if err != nil {
				return fmt.Errorf("failed to listen: %w", err)
			}

			go func() {
				log.Printf("🚀 Auth service running on port %s", grpcPort)
				if err := srv.Serve(lis); err != nil {
					log.Fatalf("server error: %v", err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Shutting down gRPC server...")
			srv.GracefulStop()
			return nil
		},
	})
}
