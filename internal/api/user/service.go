package user

import (
	desc "github.com/ArturSaga/auth/api/grpc/pkg/user_v1"
	"github.com/ArturSaga/auth/internal/service"
)

// UserAPI - сущность, которая ипмлементирует контракты
type UserAPI struct {
	desc.UnimplementedUserApiServer
	userService service.UserService
}

// NewUserAPI - публичный метод, реализует контракты
func NewUserAPI(userService service.UserService) *UserAPI {
	return &UserAPI{
		userService: userService,
	}
}
