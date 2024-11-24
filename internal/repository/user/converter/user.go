package converter

import (
	desc "github.com/ArturSaga/auth/api/grpc/pkg/user_v1"
	modelRedis "github.com/ArturSaga/auth/internal/client/cache/redis/model"
	"github.com/ArturSaga/auth/internal/constants"
	"github.com/ArturSaga/auth/internal/model"
	modelRepo "github.com/ArturSaga/auth/internal/repository/user/model"
)

// ToUserFromRepo - ковертер, который преобразует модель репо слоя в смодель сервисного слоя
func ToUserFromRepo(user *modelRepo.User) *model.User {
	if user == nil {
		return &model.User{}
	}

	return &model.User{
		ID:        user.ID,
		Info:      *ToUserInfoFromRepo(user),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// ToUserRedisFromRepo - ковертер, который преобразует модель апи (протобаф) слоя в модель сервисного слоя
func ToUserRedisFromRepo(user *modelRepo.User) *modelRedis.User {
	if user == nil {
		return &modelRedis.User{}
	}

	return &modelRedis.User{
		ID:              user.ID,
		Name:            user.Name,
		Email:           user.Email,
		Password:        user.Password,
		PasswordConfirm: "",
		Role:            user.Role,
		CreatedAt:       user.CreatedAt.UnixNano(),
		UpdatedAt:       user.UpdatedAt.UnixNano(),
	}
}

// ToUserInfoFromRepo - ковертер, который преобразует модель репо слоя в смодель сервисного слоя
func ToUserInfoFromRepo(info *modelRepo.User) *model.UserInfo {
	if info == nil {
		return &model.UserInfo{}
	}

	return &model.UserInfo{
		Name:            info.Name,
		Email:           info.Email,
		Password:        info.Password,
		PasswordConfirm: "",
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
