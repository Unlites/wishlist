package middleware

import "context"

type TokenManager interface {
	ParseUserId(ctx context.Context, token string) (int, error)
}

type MiddlewareProvider struct {
	tokenManager TokenManager
}

func NewMiddlewareProvider(tokenManager TokenManager) *MiddlewareProvider {
	return &MiddlewareProvider{
		tokenManager: tokenManager,
	}
}
