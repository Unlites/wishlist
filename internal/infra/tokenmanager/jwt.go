package tokenmanager

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTTokenManager struct {
	secretKey      string
	expirationTime time.Duration
}

func NewJWTTokenManager(secretKey string, expirationTime time.Duration) *JWTTokenManager {
	return &JWTTokenManager{
		secretKey:      secretKey,
		expirationTime: expirationTime,
	}
}

func (j *JWTTokenManager) ParseUserId(ctx context.Context, token string) (int, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Header["alg"])
		}

		return []byte(j.secretKey), nil
	})
	if err != nil {
		return 0, fmt.Errorf("jwt.Parse: %w", err)
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("failed to get user from claims")
	}

	userId, err := strconv.Atoi(claims["sub"].(string))
	if err != nil {
		return 0, fmt.Errorf("failed to get user_id from claims: %w", err)
	}

	return userId, nil
}

func (j *JWTTokenManager) GenerateToken(ctx context.Context, sub string) (string, error) {
	expiresAt := time.Now().Add(j.expirationTime).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   sub,
		ExpiresAt: expiresAt,
	})

	tokenStr, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
