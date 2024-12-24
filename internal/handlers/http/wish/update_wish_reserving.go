package wish

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Unlites/wishlist/internal/domain"
)

type updateWishReservingRequest struct {
	IsReserved bool `json:"is_reserved"`
}

func (wh *WishHandler) UpdateWishReserving(w http.ResponseWriter, r *http.Request) {
	wishId := r.PathValue("wish_id")
	wishIdInt, err := strconv.Atoi(wishId)
	if err != nil {
		http.Error(w, fmt.Errorf("strconv.Atoi: %w", err).Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	var req updateWishReservingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Errorf("json.Decode: %w", err).Error(), http.StatusBadRequest)
		return
	}

	if err := wh.service.UpdateWishReserving(ctx, wishIdInt, req.IsReserved); err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, domain.ErrNotFound) {
			status = http.StatusUnauthorized
		}

		http.Error(w, fmt.Errorf("service.UpdateWishReserving: %w", err).Error(), status)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
