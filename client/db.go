package client

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

const DSN = "postgresql://postgres:postgres@localhost:5434/postgres?sslmode=disable"

type DB struct {
	Pool *pgxpool.Pool
}

func NewDB(ctx context.Context) (*DB, error) {
	
	pool, err := pgxpool.Connect(ctx, DSN)
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
	db.Pool.Close()
}
