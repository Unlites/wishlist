package postgres

import "github.com/jackc/pgx/v5/pgxpool"

type WishRepositoryPostgres struct {
	pool *pgxpool.Pool
}

func NewWishRepositoryPostgres(pool *pgxpool.Pool) *WishRepositoryPostgres {
	return &WishRepositoryPostgres{
		pool: pool,
	}
}
