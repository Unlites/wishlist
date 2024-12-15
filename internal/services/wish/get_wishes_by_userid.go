package wish

import (
	"context"

	"github.com/Unlites/wishlist/internal/domain"
)

func (ws *WishService) GetWishesByUserId(ctx context.Context, userId int) ([]domain.Wish, error) {
	return ws.wishRepo.GetWishesByUserId(ctx, userId)
}
