package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Unlites/wishlist/internal/domain"
	"github.com/Unlites/wishlist/internal/handlers/http/response"
	validation "github.com/go-ozzo/ozzo-validation"
)

type registerRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (r *registerRequest) Validate() error {
	return validation.ValidateStruct(
		r,
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Password, validation.Required, validation.Length(4, 64)),
	)
}

type registerResponse struct {
	ID int `json:"id"`
}

func (uh *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Errorf("json.Decode: %w", err).Error(), http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, fmt.Errorf("validate: %w", err).Error(), http.StatusBadRequest)
		return
	}

	id, err := uh.service.Register(r.Context(), domain.User{
		Name:     req.Name,
		Password: req.Password,
	})
	if err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, domain.ErrNotFound) {
			status = http.StatusUnauthorized
		}

		http.Error(w, fmt.Errorf("service.Register: %w", err).Error(), status)
		return
	}

	response.JSON(w, http.StatusCreated, registerResponse{ID: id})

}
