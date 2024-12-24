package hasher

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct {
	hashCost int
}

func NewBcryptHasher(hashCost int) *BcryptHasher {
	return &BcryptHasher{
		hashCost: hashCost,
	}
}

func (h *BcryptHasher) Hash(ctx context.Context, password string) (string, error) {
	res, err := bcrypt.GenerateFromPassword([]byte(password), h.hashCost)
	return string(res), err
}

func (h *BcryptHasher) Compare(context context.Context, password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
