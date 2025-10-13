package postgres

import (
	"auth-micro/client"
	"auth-micro/internal/auth/entity"
	"auth-micro/internal/auth/repository"
	"context"
)

type userRepo struct {
	db *client.DB
}

func NewUserRepo(db *client.DB) repository.UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, u *entity.User) error {
	_, err := r.db.Pool.Exec(ctx, `
		INSERT INTO users (id, username, name, email, age, bio, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, u.ID, u.Username, u.Name, u.Email, u.Age, u.Bio, u.Password, u.CreatedAt, u.UpdatedAt)
	return err
}

func (r *userRepo) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	row := r.db.Pool.QueryRow(ctx, `
		SELECT id, username, name, email, age, bio, password, created_at, updated_at
		FROM users WHERE username = $1
	`, username)

	var u entity.User
	if err := row.Scan(&u.ID, &u.Username, &u.Name, &u.Email, &u.Age, &u.Bio, &u.Password, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}
