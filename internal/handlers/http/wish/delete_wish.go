package wish

import (
	"fmt"
	"net/http"
	"strconv"

	cctx "github.com/Unlites/wishlist/internal/common/ctx"
)

func (wh *WishHandler) DeleteWish(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("user_id")
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		http.Error(w, fmt.Errorf("strconv.Atoi: %w", err).Error(), http.StatusBadRequest)
		return
	}

	wishId := r.PathValue("wish_id")
	wishIdInt, err := strconv.Atoi(wishId)
	if err != nil {
		http.Error(w, fmt.Errorf("strconv.Atoi: %w", err).Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	callerUserId := cctx.GetUserId(ctx)
	if userIdInt != callerUserId {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	if err := wh.service.DeleteWish(ctx, wishIdInt); err != nil {
		http.Error(w, fmt.Errorf("service.DeleteWish: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
