package user

import (
	"context"
	"fmt"

	"github.com/Unlites/wishlist/internal/domain"
)

func (s *UserService) Register(ctx context.Context, user domain.User) (int, error) {
	hash, err := s.hasher.Hash(ctx, user.Password)
	if err != nil {
		return 0, fmt.Errorf("hasher.Hash: %w", err)
	}
	user.Password = hash
	return s.repo.AddUser(ctx, user)
}
