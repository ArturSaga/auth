package user

import (
	"context"
	"strconv"

	"google.golang.org/protobuf/types/known/emptypb"

	desc "github.com/ArturSaga/auth/api/grpc/pkg/user_v1"
	serviceModel "github.com/ArturSaga/auth/internal/model"
)

// UpdateUser - публичный метод, сервиса для обновления пользователя
func (s *serv) UpdateUser(ctx context.Context, userInfo *serviceModel.UpdateUserInfo) (emptypb.Empty, error) {
	if *userInfo.Role == desc.Role_UNKNOWN {
		userInfo.Role = nil
	}

	_, err := s.userRepo.UpdateUser(ctx, userInfo)
	if err != nil {
		return emptypb.Empty{}, err
	}

	err = s.cache.Del(ctx, strconv.FormatInt(userInfo.UserID, 10))
	if err != nil {
		return emptypb.Empty{}, err
	}

	return emptypb.Empty{}, nil
}
