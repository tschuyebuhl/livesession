package data

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
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
	slog.Info("getting user from postgres repository", "id", id)
	var user User
	err := r.pool.QueryRow(context.Background(), "SELECT id, name, surname FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name, &user.Surname)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Error("no user found with given id", "id", id)
			return nil, err
		} else {
			slog.Error("error getting user", "id", id, "error", err)
			return nil, err
		}
	}
	return &user, nil
}
