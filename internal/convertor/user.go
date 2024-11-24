package converter

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	desc "github.com/ArturSaga/auth/api/grpc/pkg/user_v1"
	redisModel "github.com/ArturSaga/auth/internal/client/cache/redis/model"
	"github.com/ArturSaga/auth/internal/constants"
	serviceModel "github.com/ArturSaga/auth/internal/model"
)

// ToUserDescFromService - ковертер, который преобразует модель сервисного слоя в модель апи (протобаф) слоя
func ToUserDescFromService(user *serviceModel.User) *desc.User {
	if user == nil {
		return nil
	}

	return &desc.User{
		Id:        user.ID,
		Info:      ToUserInfoDescFromService(&user.Info),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}

// ToUserInfoDescFromService - ковертер, который преобразует модель сервисного слоя в модель апи (протобаф) слоя
func ToUserInfoDescFromService(info *serviceModel.UserInfo) *desc.UserInfo {
	if info == nil {
		return nil
	}

	return &desc.UserInfo{
		Name:            info.Name,
		Email:           info.Email,
		Password:        info.Password,
		PasswordConfirm: info.PasswordConfirm,
		Role:            info.Role,
	}
}

// ToUserServiceFromRedis - ковертер, который преобразует модель апи (протобаф) слоя в модель сервисного слоя
func ToUserServiceFromRedis(user *redisModel.User) *serviceModel.User {
	if user == nil {
		return &serviceModel.User{}
	}

	userInfo := serviceModel.UserInfo{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Role:     RoleFromString(user.Role),
	}

	return &serviceModel.User{
		ID:        user.ID,
		Info:      userInfo,
		CreatedAt: time.Unix(0, user.CreatedAt),
		UpdatedAt: time.Unix(0, user.UpdatedAt),
	}
}

// ToUserRedisFromService - ковертер, который преобразует модель апи (протобаф) слоя в модель сервисного слоя
func ToUserRedisFromService(user *serviceModel.User) *redisModel.User {
	if user == nil {
		return &redisModel.User{}
	}

	return &redisModel.User{
		ID:              user.ID,
		Name:            user.Info.Name,
		Email:           user.Info.Email,
		Password:        user.Info.Password,
		PasswordConfirm: "",
		Role:            user.Info.Role.String(),
		CreatedAt:       user.CreatedAt.UnixNano(),
		UpdatedAt:       user.UpdatedAt.UnixNano(),
	}
}

// ToUserInfoServiceFromDesc - ковертер, который преобразует модель апи (протобаф) слоя в модель сервисного слоя
func ToUserInfoServiceFromDesc(info *desc.UserInfo) *serviceModel.UserInfo {
	if info == nil {
		return &serviceModel.UserInfo{}
	}

	return &serviceModel.UserInfo{
		Name:            info.Name,
		Email:           info.Email,
		Password:        info.Password,
		PasswordConfirm: info.PasswordConfirm,
		Role:            info.Role,
	}
}

// ToUpdateUserInfoServiceFromDesc - ковертер, который преобразует модель апи (протобаф) слоя в модель сервисного слоя
func ToUpdateUserInfoServiceFromDesc(info *desc.UpdateUserInfo) *serviceModel.UpdateUserInfo {
	if info == nil {
		return &serviceModel.UpdateUserInfo{}
	}

	return &serviceModel.UpdateUserInfo{
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

// RoleFromString Функция для преобразования Role из строкового представления
func RoleFromString(s string) desc.Role {
	switch s {
	case constants.ADMIN:
		return desc.Role_ADMIN
	case constants.USER:
		return desc.Role_USER
	default:
		return desc.Role_UNKNOWN
	}
}
