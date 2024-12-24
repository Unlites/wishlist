package wish

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Unlites/wishlist/internal/domain"
	"github.com/Unlites/wishlist/internal/handlers/http/response"
)

type wishResponse struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	IsReserved  bool      `json:"is_reserved"`
	UserId      int       `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
}

func (wh *WishHandler) GetWishesByUserId(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("user_id")
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		http.Error(w, fmt.Errorf("strconv.Atoi: %w", err).Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	wishes, err := wh.service.GetWishesByUserId(ctx, userIdInt)
	if err != nil {
		http.Error(w, fmt.Errorf("service.GetWishesByUserId: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	response.JSON(w, http.StatusOK, toWishResponses(wishes))
}

func toWishResponses(wishes []domain.Wish) []wishResponse {
	res := make([]wishResponse, len(wishes))
	for i, wish := range wishes {
		res[i] = wishResponse{
			Id:          wish.Id,
			Title:       wish.Title,
			Description: wish.Description,
			IsReserved:  wish.IsReserved,
			UserId:      wish.UserId,
			CreatedAt:   wish.CreatedAt,
		}
	}
	return res
}
