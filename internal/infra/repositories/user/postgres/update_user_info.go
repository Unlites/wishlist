package postgres

import (
	"context"
	"fmt"
)

func (urp *UserRepositoryPostgres) UpdateUserInfo(ctx context.Context, userId int, info string) error {
	conn, err := urp.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("pool.Acquire: %w", err)
	}
	defer conn.Release()

	query := "UPDATE wishlist.users SET info = $1 WHERE id = $2"

	if _, err := conn.Exec(ctx, query, info, userId); err != nil {
		return fmt.Errorf("conn.Exec: %w", err)
	}

	return nil
}
