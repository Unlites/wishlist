package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Unlites/wishlist/internal/domain"
	"github.com/jackc/pgx/v5"
)

func (urp *UserRepositoryPostgres) GetUserById(ctx context.Context, userId int) (domain.User, error) {
	conn, err := urp.pool.Acquire(ctx)
	if err != nil {
		return domain.User{}, fmt.Errorf("pool.Acquire: %w", err)
	}
	defer conn.Release()

	query := "SELECT id, name, password_hash FROM wishlist.users WHERE id = $1"

	var user domain.User

	if err := conn.QueryRow(ctx, query, userId).Scan(&user.Id, &user.Name, &user.Password); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrNotFound
		}

		return domain.User{}, fmt.Errorf("conn.QueryRow.Scan: %w", err)
	}

	return user, nil
}
