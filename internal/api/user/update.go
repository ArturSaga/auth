package user

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/emptypb"

	desc "github.com/ArturSaga/auth/api/grpc/pkg/user_v1"
	converter "github.com/ArturSaga/auth/internal/convertor"
)

// UpdateUser - публичный метод, который обновляет данные пользователя.
func (i *Implementation) UpdateUser(ctx context.Context, req *desc.UpdateUserRequest) (*emptypb.Empty, error) {
	_, err := i.userService.UpdateUser(ctx, converter.ToUpdateUserInfoFromDesc(req.Info))

	if err != nil {
		fmt.Printf("failed to update user: %v", err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
