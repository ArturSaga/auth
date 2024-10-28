package repository

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ArturSaga/auth/internal/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, userInfo *model.UserInfo) (int64, error)
	GetUser(ctx context.Context, id int64) (*model.User, error)
	UpdateUser(ctx context.Context, userInfo *model.UpdateUserInfo) (emptypb.Empty, error)
	DeleteUser(ctx context.Context, id int64) (emptypb.Empty, error)
}
