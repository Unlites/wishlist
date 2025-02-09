package user

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Unlites/wishlist/internal/domain"
	"github.com/Unlites/wishlist/internal/handlers/http/response"
)

type userResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (uh *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("user_id")
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		http.Error(w, fmt.Errorf("strconv.Atoi: %w", err).Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	user, err := uh.service.GetUserById(ctx, userIdInt)
	if err != nil {
		http.Error(w, fmt.Errorf("service.GetWishesByUserId: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	response.JSON(w, http.StatusOK, toUserResponse(user))
}

func toUserResponse(user domain.User) userResponse {
	return userResponse{
		Id:   user.Id,
		Name: user.Name,
	}
}
