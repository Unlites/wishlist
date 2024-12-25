package postgres

import (
	"context"
	"fmt"

	"github.com/Unlites/wishlist/internal/domain"
)

func (wrp *WishRepositoryPostgres) GetWishesByUserId(ctx context.Context, userId int) ([]domain.Wish, error) {
	conn, err := wrp.pool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("pool.Acquire: %w", err)
	}
	defer conn.Release()

	query := `
		SELECT id, title, description, is_reserved, user_id, created_at
		FROM wishlist.wishes
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	var wishes []domain.Wish

	rows, err := conn.Query(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("conn.Query: %w", err)
	}

	for rows.Next() {
		var wish domain.Wish
		if err := rows.Scan(
			&wish.Id,
			&wish.Title,
			&wish.Description,
			&wish.IsReserved,
			&wish.UserId,
			&wish.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}
		wishes = append(wishes, wish)
	}

	return wishes, nil
}
