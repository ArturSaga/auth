package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	desc "github.com/ArturSaga/auth/api/grpc/pkg/user_v1"
	"github.com/ArturSaga/auth/internal/model"
)

// ToUserFromService - ковертер, который преобразует модель сервисного слоя в модель апи (протобаф) слоя
func ToUserFromService(user *model.User) *desc.User {
	return &desc.User{
		Id:        user.ID,
		Info:      ToUserInfoFromRepo(&user.Info),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}

// ToUserInfoFromRepo - ковертер, который преобразует модель сервисного слоя в модель апи (протобаф) слоя
func ToUserInfoFromRepo(info *model.UserInfo) *desc.UserInfo {
	return &desc.UserInfo{
		Name:            info.Name,
		Email:           info.Email,
		Password:        info.Password,
		PasswordConfirm: info.PasswordConfirm,
		Role:            info.Role,
	}
}

// ToUserInfoFromDesc - ковертер, который преобразует модель апи (протобаф) слоя в модель сервисного слоя
func ToUserInfoFromDesc(info *desc.UserInfo) *model.UserInfo {
	return &model.UserInfo{
		Name:            info.Name,
		Email:           info.Email,
		Password:        info.Password,
		PasswordConfirm: info.PasswordConfirm,
		Role:            info.Role,
	}
}

// ToUpdateUserInfoFromDesc - ковертер, который преобразует модель апи (протобаф) слоя в модель сервисного слоя
func ToUpdateUserInfoFromDesc(info *desc.UpdateUserInfo) *model.UpdateUserInfo {
	return &model.UpdateUserInfo{
		UserID:          info.UserID,
		Name:            checkEmptyOrNil(info.Name),
		OldPassword:     checkEmptyOrNil(info.OldPassword),
		Password:        checkEmptyOrNil(info.Password),
		PasswordConfirm: checkEmptyOrNil(info.PasswordConfirm),
		Role:            &info.Role,
	}
}

// checkEmptyOrNil - функция, которая преобразует *wrapperspb.StringValue в *string
func checkEmptyOrNil(s *wrapperspb.StringValue) *string {
	if s == nil {
		return nil
	}

	str := s.Value
	return &str
}
