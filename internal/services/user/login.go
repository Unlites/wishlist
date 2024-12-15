package user

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Unlites/wishlist/internal/domain"
)

func (us *UserService) Login(ctx context.Context, user domain.User) (string, error) {
	existingUser, err := us.repo.GetUserById(ctx, user.Id)
	if err != nil {
		return "", fmt.Errorf("GetUserById: %w", err)
	}

	if !us.hasher.Compare(ctx, existingUser.Password, user.Password) {
		return "", fmt.Errorf("%v, wrong password", domain.ErrUnauthorized)
	}

	token, err := us.tokenManager.GenerateToken(ctx, strconv.Itoa(user.Id))
	if err != nil {
		return "", fmt.Errorf("GenerateToken: %w", err)
	}

	return token, nil
}
