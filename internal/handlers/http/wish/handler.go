package wish

import (
	"context"
	"net/http"

	"github.com/Unlites/wishlist/internal/domain"
	"github.com/Unlites/wishlist/internal/handlers/http/middleware"
)

type WishService interface {
	GetWishesByUserId(ctx context.Context, userId int) ([]domain.Wish, error)
	AddWish(ctx context.Context, wish domain.Wish) (int, error)
	UpdateWish(ctx context.Context, wish domain.Wish) error
	DeleteWish(ctx context.Context, wishId int) error
	UpdateWishReserving(ctx context.Context, wishId int, isReserved bool) error
}

type WishHandler struct {
	service    WishService
	mwProvider *middleware.MiddlewareProvider
}

func NewWishHandler(wishService WishService, middlewareProvider *middleware.MiddlewareProvider) *WishHandler {
	return &WishHandler{
		service:    wishService,
		mwProvider: middlewareProvider,
	}
}

func (wh *WishHandler) RegisterRoutes(mux *http.ServeMux, prefix string) {
	mux.Handle("POST "+prefix, wh.mwProvider.Auth(http.HandlerFunc(wh.AddWish)))
	mux.Handle("PUT "+prefix+"/{wish_id}", wh.mwProvider.Auth(http.HandlerFunc(wh.UpdateWish)))
	mux.Handle("DELETE "+prefix+"/{wish_id}", wh.mwProvider.Auth(http.HandlerFunc(wh.DeleteWish)))
	mux.Handle("PUT "+prefix+"/{wish_id}/update-reserving", wh.mwProvider.Auth(http.HandlerFunc(wh.UpdateWishReserving)))
	mux.Handle("GET "+prefix, wh.mwProvider.Auth(http.HandlerFunc(wh.GetWishesByUserId)))
}
