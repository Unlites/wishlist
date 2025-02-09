package user

import (
	"context"

	"github.com/Unlites/wishlist/internal/domain"
)

func (us *UserService) GetUserById(ctx context.Context, userId int) (domain.User, error) {
	return us.repo.GetUserById(ctx, userId)
}
