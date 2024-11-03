package user

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/emptypb"

	desc "github.com/ArturSaga/auth/api/grpc/pkg/user_v1"
)

// DeleteUser - публичный метод, который удаляет пользователя.
func (i *UserAPI) DeleteUser(ctx context.Context, req *desc.DeleteUserRequest) (*emptypb.Empty, error) {
	_, err := i.userService.DeleteUser(ctx, req.Id)
	if err != nil {
		fmt.Printf("failed to delete user: %v", err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
