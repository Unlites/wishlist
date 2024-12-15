package wish

import (
	"context"
	"fmt"
)

func (ws *WishService) UpdateWishReserving(ctx context.Context, wishId int, isReserved bool) error {
	wish, err := ws.wishRepo.GetWishById(ctx, wishId)
	if err != nil {
		return fmt.Errorf("wishRepository.GetWishById: %w", err)
	}

	wish.IsReserved = isReserved
	return ws.wishRepo.UpdateWish(ctx, wish)
}
