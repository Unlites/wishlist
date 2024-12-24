package postgres

import (
	"context"
	"fmt"
)

func (wrp *WishRepositoryPostgres) DeleteWish(ctx context.Context, wishId int) error {
	conn, err := wrp.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("pool.Acquire: %w", err)
	}
	defer conn.Release()

	query := "DELETE FROM wishlist.wishes WHERE id = $1"

	if _, err := conn.Exec(ctx, query, wishId); err != nil {
		return fmt.Errorf("conn.Exec: %w", err)
	}

	return nil
}
