package user

import (
	"context"
	"fmt"

	desc "github.com/ArturSaga/auth/api/grpc/pkg/user_v1"
	converter "github.com/ArturSaga/auth/internal/convertor"
)

// CreateUser - публичный метод, который создает пользователя.
func (i *Implementation) CreateUser(ctx context.Context, req *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {
	id, err := i.userService.CreateUser(ctx, converter.ToUserInfoFromDesc(req.Info))

	if err != nil {
		fmt.Printf("failed to create user: %v", err)
		return nil, err
	}

	return &desc.CreateUserResponse{
		Id: id,
	}, nil
}
