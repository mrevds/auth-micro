package middleware

import (
	"context"

	"go.uber.org/ratelimit"
	"google.golang.org/grpc"
)

func RateLimitInterceptor(rl ratelimit.Limiter) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		rl.Take()
		return handler(ctx, req)
	}
}
