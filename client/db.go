package client

import (
    "context"
    "fmt"
    "os"
    "time"

    "github.com/jackc/pgx/v4/pgxpool"
    "go.uber.org/fx"
)

type DB struct {
    Pool *pgxpool.Pool
}

// Params для dependency injection
type Params struct {
    fx.In
    Lifecycle fx.Lifecycle
}

// NewDB создает подключение к БД с lifecycle hooks
func NewDB(p Params) (*DB, error) {
    // Построение DSN из переменных окружения
    dsn := fmt.Sprintf(
        "postgresql://%s:%s@%s:%s/%s?sslmode=disable",
        getEnv("PG_USER", "auth_db_user"),
        getEnv("PG_PASSWORD", "auth_db_password"),
        getEnv("PG_HOST", "localhost"),
        getEnv("PG_PORT", "54322"),
        getEnv("PG_DATABASE_NAME", "auth_db"),
    )

    // Парсим конфигурацию
    config, err := pgxpool.ParseConfig(dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to parse DSN: %w", err)
    }

    // Настройка пула
    config.MaxConns = 25
    config.MinConns = 5
    config.MaxConnLifetime = time.Hour
    config.MaxConnIdleTime = 30 * time.Minute

    var pool *pgxpool.Pool
    db := &DB{}

    // Lifecycle hooks
    p.Lifecycle.Append(fx.Hook{
        OnStart: func(ctx context.Context) error {
            pool, err = pgxpool.ConnectConfig(ctx, config)
            if err != nil {
                return fmt.Errorf("failed to connect to database: %w", err)
            }

            if err = pool.Ping(ctx); err != nil {
                pool.Close()
                return fmt.Errorf("failed to ping database: %w", err)
            }

            db.Pool = pool
            fmt.Println("✅ Database connected successfully")
            return nil
        },
        OnStop: func(ctx context.Context) error {
            if db.Pool != nil {
                db.Pool.Close()
                fmt.Println("✅ Database connection closed")
            }
            return nil
        },
    })

    return db, nil
}

// getEnv получает значение переменной окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}