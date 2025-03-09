package user

import "context"

func (s *UserService) UpdateUserInfo(ctx context.Context, userId int, info string) error {
	return s.repo.UpdateUserInfo(ctx, userId, info)
}
