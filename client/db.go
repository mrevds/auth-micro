package client

import (
    "context"
    "fmt"
    "os"

    "github.com/jackc/pgx/v4/pgxpool"
)

type DB struct {
    Pool *pgxpool.Pool
}

func NewDB(ctx context.Context) (*DB, error) {
    // Построение DSN из переменных окружения
    dsn := fmt.Sprintf(
        "postgresql://%s:%s@%s:%s/%s?sslmode=disable",
        getEnv("PG_USER", "auth_db_user"),
        getEnv("PG_PASSWORD", "auth_db_password"),
        getEnv("PG_HOST", "localhost"),
        getEnv("PG_PORT", "54322"),
        getEnv("PG_DATABASE_NAME", "auth_db"),
    )
    
    pool, err := pgxpool.Connect(ctx, dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }

    if err = pool.Ping(ctx); err != nil {
        pool.Close()
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }

    return &DB{Pool: pool}, nil
}

func (db *DB) Close() {
    if db.Pool != nil {
        db.Pool.Close()
    }
}

// getEnv получает значение переменной окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}