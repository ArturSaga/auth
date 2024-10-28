package user

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ArturSaga/auth/internal/model"
)

func (s *serv) UpdateUser(ctx context.Context, userInfo *model.UpdateUserInfo) (emptypb.Empty, error) {
	_, err := s.userRepo.UpdateUser(ctx, userInfo)
	if err != nil {
		return emptypb.Empty{}, err
	}

	return emptypb.Empty{}, nil
}
