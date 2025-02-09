package user

import (
	"context"
	"net/http"

	"github.com/Unlites/wishlist/internal/domain"
)

type UserService interface {
	Login(ctx context.Context, user domain.User) (string, error)
	Register(ctx context.Context, user domain.User) (int, error)
	GetUserById(ctx context.Context, userId int) (domain.User, error)
}

type UserHandler struct {
	service UserService
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (uh *UserHandler) RegisterRoutes(mux *http.ServeMux, prefix string) {
	mux.Handle("POST "+prefix+"/login", http.HandlerFunc(uh.Login))
	mux.Handle("POST "+prefix+"/register", http.HandlerFunc(uh.Register))
	mux.Handle("GET "+prefix+"/{user_id}", http.HandlerFunc(uh.GetUserById))
}
