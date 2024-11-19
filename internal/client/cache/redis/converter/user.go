package converter

import (
	desc "github.com/ArturSaga/auth/api/grpc/pkg/user_v1"
	modelRedis "github.com/ArturSaga/auth/internal/client/cache/redis/model"
	"github.com/ArturSaga/auth/internal/constants"
	"github.com/ArturSaga/auth/internal/model"
)

// ToUserInfoFromRedis - ковертер, который преобразует модель репо слоя в смодель сервисного слоя
func ToUserInfoFromRedis(info *modelRedis.User) *model.UserInfo {
	if info == nil {
		return &model.UserInfo{}
	}

	return &model.UserInfo{
		Name:            info.Name,
		Email:           info.Email,
		Password:        info.Password,
		PasswordConfirm: info.PasswordConfirm,
		Role:            RoleFromString(info.Role),
	}
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
