package wish

import (
	"context"
	"fmt"

	cctx "github.com/Unlites/wishlist/internal/common/ctx"
	"github.com/Unlites/wishlist/internal/domain"
)

func (ws *WishService) GetWishesByUserId(ctx context.Context, userId int) ([]domain.Wish, error) {
	wishes, err := ws.wishRepo.GetWishesByUserId(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("wishRepository.GetWishesByUserId: %w", err)
	}

	for i := range wishes {
		if cctx.GetUserId(ctx) == userId {
			// do not show isReserved for owner
			wishes[i].IsReserved = nil
		}
	}

	return wishes, nil
}
