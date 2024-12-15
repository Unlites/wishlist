package user

import (
	"context"

	"github.com/Unlites/wishlist/internal/domain"
)

type UserService interface {
	Login(ctx context.Context, user domain.User) (string, error)
	Register(ctx context.Context, user domain.User) (int, error)
}
