package user

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *serv) DeleteUser(ctx context.Context, id int64) (emptypb.Empty, error) {
	_, err := s.userRepo.DeleteUser(ctx, id)
	if err != nil {
		return emptypb.Empty{}, err
	}
	return emptypb.Empty{}, nil
}
