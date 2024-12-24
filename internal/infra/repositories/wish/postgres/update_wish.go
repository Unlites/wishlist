package postgres

import (
	"context"
	"fmt"

	"github.com/Unlites/wishlist/internal/domain"
)

func (wrp *WishRepositoryPostgres) UpdateWish(ctx context.Context, wish domain.Wish) error {
	conn, err := wrp.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("pool.Acquire: %w", err)
	}
	defer conn.Release()

	query := `
		UPDATE wishlist.wishes
		SET title = $1, description = $2, is_reserved = $3
		WHERE id = $4
	`

	if _, err := conn.Exec(ctx,
		query,
		wish.Title,
		wish.Description,
		wish.IsReserved,
		wish.Id,
	); err != nil {
		return fmt.Errorf("conn.Exec: %w", err)
	}

	return nil
}
