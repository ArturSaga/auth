package user

import (
	"context"
	"errors"
	"fmt"

	desc "github.com/ArturSaga/auth/api/grpc/pkg/user_v1"
	converter "github.com/ArturSaga/auth/internal/convertor"
)

// CreateUser - публичный метод, который создает пользователя.
func (i *Implementation) CreateUser(ctx context.Context, req *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {
	if req.Info.Password != "" && req.Info.PasswordConfirm != "" {
		if req.Info.Password != req.Info.PasswordConfirm {
			return nil, errors.New("passwords don't match")
		}
	} else {
		return nil, errors.New("passwords can't be empty")
	}
	if req.Info.Name == "" {
		return nil, errors.New("name can't be empty")
	}
	if req.Info.Email == "" {
		return nil, errors.New("email can't be empty")
	}

	id, err := i.userService.CreateUser(ctx, converter.ToUserInfoFromDesc(req.Info))

	if err != nil {
		fmt.Printf("failed to create user: %v", err)
		return nil, err
	}

	return &desc.CreateUserResponse{
		Id: id,
	}, nil
}
