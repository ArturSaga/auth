package converter

import (
	desc "github.com/ArturSaga/auth/api/grpc/pkg/user_v1"
	"github.com/ArturSaga/auth/internal/model"
	modelRepo "github.com/ArturSaga/auth/internal/repository/user/model"
)

// ToUserFromRepo - ковертер, который преобразует модель репо слоя в смодель сервисного слоя
func ToUserFromRepo(user *modelRepo.User) *model.User {
	return &model.User{
		ID:        user.ID,
		Info:      ToUserInfoFromRepo(user),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// ToUserInfoFromRepo - ковертер, который преобразует модель репо слоя в смодель сервисного слоя
func ToUserInfoFromRepo(info *modelRepo.User) model.UserInfo {
	return model.UserInfo{
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
	case "UNKNOWN":
		return desc.Role_UNKNOWN
	case "ADMIN":
		return desc.Role_ADMIN
	case "USER":
		return desc.Role_USER
	default:
		return desc.Role_UNKNOWN
	}
}
