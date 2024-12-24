package wish

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	cctx "github.com/Unlites/wishlist/internal/common/ctx"
	"github.com/Unlites/wishlist/internal/domain"
	"github.com/Unlites/wishlist/internal/handlers/http/response"
	validation "github.com/go-ozzo/ozzo-validation"
)

type addWishRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (r *addWishRequest) Validate() error {
	return validation.ValidateStruct(
		r,
		validation.Field(&r.Title, validation.Required, validation.Length(1, 100)),
		validation.Field(&r.Description, validation.Length(1, 1000)),
	)
}

type addWishResponse struct {
	ID int `json:"id"`
}

func (wh *WishHandler) AddWish(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("user_id")
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		http.Error(w, fmt.Errorf("strconv.Atoi: %w", err).Error(), http.StatusBadRequest)
		return
	}

	var req addWishRequest
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

	id, err := wh.service.AddWish(ctx, domain.Wish{
		Title:       req.Title,
		Description: req.Description,
		IsReserved:  false,
		UserId:      userIdInt,
	})
	if err != nil {
		http.Error(w, fmt.Errorf("service.AddWish: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	response.JSON(w, http.StatusCreated, addWishResponse{ID: id})
}
