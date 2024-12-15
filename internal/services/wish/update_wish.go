package wish

import (
	"context"

	"github.com/Unlites/wishlist/internal/domain"
)

func (ws *WishService) UpdateWish(ctx context.Context, wish domain.Wish) error {
	return ws.wishRepo.UpdateWish(ctx, wish)
}
