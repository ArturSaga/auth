package user

import (
	"context"
	"fmt"

	desc "github.com/ArturSaga/auth/api/grpc/pkg/user_v1"
	converter "github.com/ArturSaga/auth/internal/convertor"
)

// GetUser - публичный метод, который позволяет получить данные пользователя.
func (i *UserApi) GetUser(ctx context.Context, req *desc.GetUserRequest) (*desc.GetUserResponse, error) {
	userObj, err := i.userService.GetUser(ctx, req.GetId())
	if err != nil {
		fmt.Printf("failed to get user: %v", err)
		return nil, err
	}

	return &desc.GetUserResponse{
		User: converter.ToUserDescFromService(userObj),
	}, nil
}
