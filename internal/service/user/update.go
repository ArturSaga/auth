package user

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	desc "github.com/ArturSaga/auth/api/grpc/pkg/user_v1"
	"github.com/ArturSaga/auth/internal/model"
)

func (s *serv) UpdateUser(ctx context.Context, userInfo *model.UpdateUserInfo) (emptypb.Empty, error) {
	if *userInfo.Role == desc.Role_UNKNOWN {
		userInfo.Role = nil
	}

	_, err := s.userRepo.UpdateUser(ctx, userInfo)
	if err != nil {
		return emptypb.Empty{}, err
	}

	return emptypb.Empty{}, nil
}
