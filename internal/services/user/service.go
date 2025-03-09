package user

import (
	"context"

	"github.com/Unlites/wishlist/internal/domain"
)

type UserRepository interface {
	GetUserById(ctx context.Context, userId int) (domain.User, error)
	GetUserByName(ctx context.Context, name string) (domain.User, error)
	AddUser(ctx context.Context, user domain.User) (int, error)
	UpdateUserInfo(ctx context.Context, userId int, info string) error
}

type Hasher interface {
	Hash(ctx context.Context, password string) (string, error)
	Compare(context context.Context, password, hash string) bool
}

type TokenManager interface {
	GenerateToken(ctx context.Context, sub string) (string, error)
}

type UserService struct {
	repo         UserRepository
	hasher       Hasher
	tokenManager TokenManager
}

func NewUserService(
	repo UserRepository,
	hasher Hasher,
	tokenManager TokenManager,
) *UserService {
	return &UserService{
		repo:         repo,
		hasher:       hasher,
		tokenManager: tokenManager,
	}

}
