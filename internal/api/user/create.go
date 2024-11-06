package user

import (
	"context"
	"fmt"
	"log"

	desc "github.com/ArturSaga/auth/api/grpc/pkg/user_v1"
	converter "github.com/ArturSaga/auth/internal/convertor"
	serviceErr "github.com/ArturSaga/auth/internal/service_error"
)

// CreateUser - публичный метод, который создает пользователя.
func (i *UserApi) CreateUser(ctx context.Context, req *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {
	if !i.validate(req) {
		return nil, serviceErr.ErrPasswordsNotMatch
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

// validate - приватный метод, проверяющий на валидность входящие данные
func (i *UserApi) validate(req *desc.CreateUserRequest) bool {
	if req.Info.Password != "" && req.Info.PasswordConfirm != "" {
		if req.Info.Password != req.Info.PasswordConfirm {
			log.Println("passwords don't match")
			return false
		}
	} else {
		log.Println("passwords can't be empty")
		return false
	}
	if req.Info.Name == "" {
		log.Println("name can't be empty")
		return false
	}
	if req.Info.Email == "" {
		log.Println("email can't be empty")
		return false
	}

	return true
}
