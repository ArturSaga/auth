package converter

import (
	"github.com/ArturSaga/auth/internal/model"
	modelRepo "github.com/ArturSaga/auth/internal/repository/user/model"
)

// ToUserFromRepo - ковертер, который преобразует модель репо слоя в смодель сервисного слоя
func ToUserFromRepo(user *modelRepo.User) *model.User {
	return &model.User{
		ID:        user.ID,
		Info:      ToUserInfoFromRepo(user.Info),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// ToUserInfoFromRepo - ковертер, который преобразует модель репо слоя в смодель сервисного слоя
func ToUserInfoFromRepo(info modelRepo.UserInfo) model.UserInfo {
	return model.UserInfo{
		Name:            info.Name,
		Email:           info.Email,
		Password:        info.Password,
		PasswordConfirm: info.PasswordConfirm,
		Role:            info.Role,
	}
}
