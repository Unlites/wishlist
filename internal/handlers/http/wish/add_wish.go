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
	Price       *int   `json:"price"`
}

func (r *addWishRequest) Validate() error {
	return validation.ValidateStruct(
		r,
		validation.Field(&r.Title, validation.Required, validation.Length(1, 300)),
		validation.Field(&r.Description, validation.Length(1, 5000)),
		validation.Field(&r.Price, validation.Min(1)),
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
		Price:       req.Price,
		UserId:      userIdInt,
	})
	if err != nil {
		http.Error(w, fmt.Errorf("service.AddWish: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	response.JSON(w, http.StatusCreated, addWishResponse{ID: id})
}
