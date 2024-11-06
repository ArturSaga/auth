package user

import (
	desc "github.com/ArturSaga/auth/api/grpc/pkg/user_v1"
	"github.com/ArturSaga/auth/internal/service"
)

// UserApi - сущность, которая ипмлементирует контракты
type UserApi struct {
	desc.UnimplementedUserApiServer
	userService service.UserService
}

// NewUserAPI - публичный метод, реализует контракты
func NewUserAPI(userService service.UserService) *UserApi {
	return &UserApi{
		userService: userService,
	}
}
