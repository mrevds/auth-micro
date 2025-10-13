package handler

import (
	"context"

	auth "auth-micro/pkg/auth_v1"
)

// AuthHandler — интерфейс хэндлера уровня транспорта.
// Обычно совпадает с protobuf-интерфейсом (gRPC) или с HTTP-роутами.
type AuthHandler interface {
	auth.AuthServer // наследуем gRPC интерфейс

	// Дополнительно можно объявить хэндлеры под другие типы транспорта.
	// Например:
	// RegisterHTTP(router chi.Router)
	// ConsumeKafkaMessage(ctx context.Context, msg []byte) error
}

// Здесь можно добавить отдельный контекстный метод для graceful shutdown или healthcheck.
type HealthCheck interface {
	Ping(ctx context.Context) error
}
