package wish

import (
	"context"

	"github.com/Unlites/wishlist/internal/domain"
)

func (ws *WishService) AddWish(ctx context.Context, wish domain.Wish) (int, error) {
	return ws.wishRepo.AddWish(ctx, wish)
}
