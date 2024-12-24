package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	cctx "github.com/Unlites/wishlist/internal/common/ctx"
)

func (mp *MiddlewareProvider) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "authorization header is empty", http.StatusUnauthorized)
			return
		}

		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		if len(authHeaderParts[1]) == 0 {
			http.Error(w, "access token is empty", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()

		userId, err := mp.tokenManager.ParseUserId(ctx, authHeaderParts[1])
		if err != nil {
			http.Error(w, fmt.Errorf("tokenManager.ParseUserId: %w", err).Error(), http.StatusUnauthorized)
			return
		}

		ctx = context.WithValue(ctx, cctx.UserIdCtxKey, userId)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
