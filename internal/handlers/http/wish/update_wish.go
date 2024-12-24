package wish

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	cctx "github.com/Unlites/wishlist/internal/common/ctx"
	"github.com/Unlites/wishlist/internal/domain"
	validation "github.com/go-ozzo/ozzo-validation"
)

type updateWishRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	IsReserved  bool   `json:"is_reserved"`
}

func (r *updateWishRequest) Validate() error {
	return validation.ValidateStruct(
		r,
		validation.Field(&r.Title, validation.Required, validation.Length(1, 100)),
		validation.Field(&r.Description, validation.Length(1, 1000)),
		validation.Field(&r.IsReserved, validation.Required),
	)
}

func (wh *WishHandler) UpdateWish(w http.ResponseWriter, r *http.Request) {
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

	var req updateWishRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Errorf("json.Decode: %w", err).Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	callerUserId := cctx.GetUserId(ctx)
	if userIdInt != callerUserId {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, fmt.Errorf("validate: %w", err).Error(), http.StatusBadRequest)
		return
	}

	if err := wh.service.UpdateWish(ctx, domain.Wish{
		Id:          wishIdInt,
		Title:       req.Title,
		Description: req.Description,
		IsReserved:  req.IsReserved,
		UserId:      userIdInt,
	}); err != nil {
		http.Error(w, fmt.Errorf("service.UpdateWish: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
