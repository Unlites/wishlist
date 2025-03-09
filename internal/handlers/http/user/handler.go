package user

import (
	"context"
	"net/http"

	"github.com/Unlites/wishlist/internal/domain"
	"github.com/Unlites/wishlist/internal/handlers/http/middleware"
)

type UserService interface {
	Login(ctx context.Context, user domain.User) (string, error)
	Register(ctx context.Context, user domain.User) (int, error)
	GetUserById(ctx context.Context, userId int) (domain.User, error)
	UpdateUserInfo(ctx context.Context, userId int, info string) error
}

type UserHandler struct {
	service    UserService
	mwProvider *middleware.MiddlewareProvider
}

func NewUserHandler(service UserService, middlewareProvider *middleware.MiddlewareProvider) *UserHandler {
	return &UserHandler{
		service:    service,
		mwProvider: middlewareProvider,
	}
}

func (uh *UserHandler) RegisterRoutes(mux *http.ServeMux, prefix string) {
	mux.Handle("POST "+prefix+"/login", http.HandlerFunc(uh.Login))
	mux.Handle("POST "+prefix+"/register", http.HandlerFunc(uh.Register))
	mux.Handle("GET "+prefix+"/{user_id}", http.HandlerFunc(uh.GetUserById))
	mux.Handle("PUT "+prefix+"/{user_id}/info", uh.mwProvider.Auth(http.HandlerFunc(uh.UpdateUserInfo)))
}
