package user

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/emptypb"

	desc "github.com/ArturSaga/auth/api/grpc/pkg/user_v1"
	converter "github.com/ArturSaga/auth/internal/convertor"
	serviceErr "github.com/ArturSaga/auth/internal/service_error"
)

// UpdateUser - публичный метод, который обновляет данные пользователя.
func (i *UserAPI) UpdateUser(ctx context.Context, req *desc.UpdateUserRequest) (*emptypb.Empty, error) {
	user, err := i.userService.GetUser(ctx, req.Info.UserID)
	updateUserInfo := converter.ToUpdateUserInfoFromDesc(req.Info)
	if err != nil {
		return &emptypb.Empty{}, serviceErr.ErrUpdateUser
	}

	if updateUserInfo.OldPassword != nil {
		err = bcrypt.CompareHashAndPassword([]byte(user.Info.Password), []byte(*updateUserInfo.OldPassword))
		if err != nil {
			return &emptypb.Empty{}, fmt.Errorf("old password not equal to current password: %v", err)
		}

		if updateUserInfo.Password != updateUserInfo.PasswordConfirm {
			return &emptypb.Empty{}, serviceErr.ErrPasswordsNotMatch
		}
	}

	_, err = i.userService.UpdateUser(ctx, updateUserInfo)

	if err != nil {
		fmt.Printf("failed to update user: %v", err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
