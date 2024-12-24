package postgres

import (
	"context"
	"fmt"

	"github.com/Unlites/wishlist/internal/domain"
	"github.com/jackc/pgx/v5/pgconn"
)

func (urp *UserRepositoryPostgres) AddUser(ctx context.Context, user domain.User) (int, error) {
	conn, err := urp.pool.Acquire(ctx)
	if err != nil {
		return 0, fmt.Errorf("pool.Acquire: %w", err)
	}
	defer conn.Release()

	query := `
		INSERT INTO wishlist.users (name, password_hash) 
		VALUES ($1, $2) RETURNING id
	`

	if err := conn.QueryRow(ctx, query, user.Name, user.Password).Scan(&user.Id); err != nil {
		pgErr, ok := err.(*pgconn.PgError)
		if ok && pgErr.Code == PgDuplicateErrorCode {
			return 0, fmt.Errorf("user with name %s %w", user.Name, domain.ErrAlreadyExists)
		}
		return 0, fmt.Errorf("conn.QueryRow.Scan: %w", err)
	}

	return user.Id, nil
}
