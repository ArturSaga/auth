package repository

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	serviceModel "github.com/ArturSaga/auth/internal/model"
)

// UserRepository - интерфейс, определящий методы репо слоя
type UserRepository interface {
	CreateUser(ctx context.Context, userInfo *serviceModel.UserInfo) (int64, error)
	GetUser(ctx context.Context, id int64) (*serviceModel.User, error)
	UpdateUser(ctx context.Context, userInfo *serviceModel.UpdateUserInfo) (emptypb.Empty, error)
	DeleteUser(ctx context.Context, id int64) (emptypb.Empty, error)
}
