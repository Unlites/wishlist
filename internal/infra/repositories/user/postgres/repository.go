package postgres

import "github.com/jackc/pgx/v5/pgxpool"

type UserRepositoryPostgres struct {
	pool *pgxpool.Pool
}

func NewUserRepositoryPostgres(pool *pgxpool.Pool) *UserRepositoryPostgres {
	return &UserRepositoryPostgres{
		pool: pool,
	}
}
