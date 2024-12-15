package wish

import (
	"context"

	"github.com/Unlites/wishlist/internal/domain"
)

type WishService interface {
	GetWishesByUserId(ctx context.Context, userId int) ([]domain.Wish, error)
	AddWish(ctx context.Context, wish domain.Wish) (int, error)
	UpdateWish(ctx context.Context, wish domain.Wish) error
	DeleteWish(ctx context.Context, wishId int) error
	UpdateWishReserving(ctx context.Context, wishId int, isReserved bool) error
}

type WishHandler struct {
	WishService WishService
}

func NewWishHandler(wishService WishService) *WishHandler {
	return &WishHandler{
		WishService: wishService,
	}
}
