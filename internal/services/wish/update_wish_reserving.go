package wish

import (
	"context"
	"fmt"

	cctx "github.com/Unlites/wishlist/internal/common/ctx"
	"github.com/Unlites/wishlist/internal/domain"
)

func (ws *WishService) UpdateWishReserving(ctx context.Context, wishId int, isReserved bool) error {
	wish, err := ws.wishRepo.GetWishById(ctx, wishId)
	if err != nil {
		return fmt.Errorf("wishRepository.GetWishById: %w", err)
	}

	reservingUserId := cctx.GetUserId(ctx)

	if wish.IsReserved != nil {
		if *wish.IsReserved == isReserved {
			return domain.ErrAlreadyProcessed
		}

		if !isReserved && wish.ReservedBy != nil && *wish.ReservedBy != reservingUserId {
			return domain.ErrForbidden
		}
	}

	wish.SetReserved(isReserved, reservingUserId)

	return ws.wishRepo.UpdateWishReserving(ctx, wish)
}
