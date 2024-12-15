package wish

import (
	"context"

	"github.com/Unlites/wishlist/internal/domain"
)

type WishRepository interface {
	GetWishesByUserId(ctx context.Context, userId int) ([]domain.Wish, error)
	GetWishById(ctx context.Context, wishId int) (domain.Wish, error)
	AddWish(ctx context.Context, wish domain.Wish) (int, error)
	UpdateWish(ctx context.Context, wish domain.Wish) error
	DeleteWish(ctx context.Context, wishId int) error
}

type WishService struct {
	wishRepo WishRepository
}

func NewWishService(wishRepository WishRepository) *WishService {
	return &WishService{
		wishRepo: wishRepository,
	}
}
