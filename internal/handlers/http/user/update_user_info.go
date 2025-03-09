package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	cctx "github.com/Unlites/wishlist/internal/common/ctx"
	validation "github.com/go-ozzo/ozzo-validation"
)

type updateUserInfoRequest struct {
	Info string `json:"info"`
}

func (r *updateUserInfoRequest) Validate() error {
	return validation.ValidateStruct(
		r, validation.Field(&r.Info, validation.Length(0, 3000)),
	)
}

func (uh *UserHandler) UpdateUserInfo(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("user_id")
	userIdInt, err := strconv.Atoi(userId)
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

	var req updateUserInfoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Errorf("json.Decode: %w", err).Error(), http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, fmt.Errorf("validate: %w", err).Error(), http.StatusBadRequest)
		return
	}

	if err := uh.service.UpdateUserInfo(ctx, userIdInt, req.Info); err != nil {
		http.Error(w, fmt.Errorf("service.UpdateUserInfo: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
