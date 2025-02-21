package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Unlites/wishlist/internal/domain"
	"github.com/jackc/pgx"
)

func (wrp *WishRepositoryPostgres) GetWishById(ctx context.Context, wishId int) (domain.Wish, error) {
	conn, err := wrp.pool.Acquire(ctx)
	if err != nil {
		return domain.Wish{}, fmt.Errorf("pool.Acquire: %w", err)
	}
	defer conn.Release()

	query := `
		SELECT id, title, description, is_reserved, reserved_by, user_id, created_at
		FROM wishlist.wishes
		WHERE id = $1
	`

	var wish domain.Wish

	if err := conn.QueryRow(ctx, query, wishId).Scan(
		&wish.Id,
		&wish.Title,
		&wish.Description,
		&wish.IsReserved,
		&wish.ReservedBy,
		&wish.UserId,
		&wish.CreatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Wish{}, domain.ErrNotFound
		}

		return domain.Wish{}, fmt.Errorf("conn.QueryRow.Scan: %w", err)
	}

	return wish, nil
}
