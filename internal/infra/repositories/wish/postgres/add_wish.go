package postgres

import (
	"context"
	"fmt"

	"github.com/Unlites/wishlist/internal/domain"
)

func (wrp *WishRepositoryPostgres) AddWish(ctx context.Context, wish domain.Wish) (int, error) {
	conn, err := wrp.pool.Acquire(ctx)
	if err != nil {
		return 0, fmt.Errorf("pool.Acquire: %w", err)
	}
	defer conn.Release()

	query := `
		INSERT INTO wishlist.wishes (title, description, is_reserved, user_id, created_at) 
		VALUES ($1, $2, $3, $4, NOW()) RETURNING id
	`

	if err := conn.QueryRow(
		ctx,
		query,
		wish.Title,
		wish.Description,
		wish.IsReserved,
		wish.UserId,
	).Scan(&wish.Id); err != nil {
		return 0, fmt.Errorf("conn.QueryRow.Scan: %w", err)
	}

	return wish.Id, nil
}
