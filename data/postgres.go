package data

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		pool: pool,
	}
}

func (r *PostgresRepository) Get(id ID) (*User, error) {
	var user *User
	err := r.pool.QueryRow(context.Background(), "SELECT id, name, surname FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name, &user.Surname)
	if err != nil {
		slog.Error("error getting user", "id", id, "error", err)
		return nil, err
	}
	return user, err
}
