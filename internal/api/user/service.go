package user

import (
	desc "github.com/ArturSaga/auth/api/grpc/pkg/user_v1"
	"github.com/ArturSaga/auth/internal/service"
)

// Implementation - сущность, которая ипмлементирует контракты
type Implementation struct {
	desc.UnimplementedUserApiServer
	userService service.UserService
}

// NewImplementation - публичный метод, реализует контракты
func NewImplementation(userService service.UserService) *Implementation {
	return &Implementation{
		userService: userService,
	}
}
