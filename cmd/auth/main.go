package main

import (
	"auth-micro/client"
	authmod "auth-micro/internal/auth"
	auth "auth-micro/pkg/auth_v1"
	"context"
	"fmt"
	"log"
	"net"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = "50051"

func main() {
	_ = godotenv.Load()

	ctx := context.Background()
	db, err := client.NewDB(ctx)
	if err != nil {
		log.Fatalf("db connection failed: %v", err)
	}

	authModule := authmod.NewModule(db)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	auth.RegisterAuthServer(srv, authModule.Handler)
	reflection.Register(srv)

	log.Printf("🚀 Auth service running on port %s", grpcPort)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
