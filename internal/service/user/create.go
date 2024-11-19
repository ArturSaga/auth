package user

import (
	"context"

	serviceModel "github.com/ArturSaga/auth/internal/model"
)

// CreateUser - публичный метод, сервиса для создания пользователя
func (s *serv) CreateUser(ctx context.Context, userInfo *serviceModel.UserInfo) (int64, error) {
	id, err := s.userRepo.CreateUser(ctx, userInfo)
	if err != nil {
		return 0, err
	}

	return id, nil
}
