package user

import (
	"context"

	"github.com/ArturSaga/auth/internal/model"
)

// GetUser - публичный метод, сервиса для получения пользователя
func (s *serv) GetUser(ctx context.Context, id int64) (*model.User, error) {
	user, err := s.userRepo.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
