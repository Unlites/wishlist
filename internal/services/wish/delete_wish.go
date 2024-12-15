package wish

import "context"

func (ws *WishService) DeleteWish(ctx context.Context, wishId int) error {
	return ws.wishRepo.DeleteWish(ctx, wishId)
}
