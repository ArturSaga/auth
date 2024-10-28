package user

import (
	"github.com/ArturSaga/auth/internal/client/db"
	"github.com/ArturSaga/auth/internal/repository"
	"github.com/ArturSaga/auth/internal/service"
)

type serv struct {
	userRepo  repository.UserRepository
	txManager db.TxManager
}

// NewUserService - публчиный метод, создающий сущность, для работы с сервисным слоем
func NewUserService(userRepo repository.UserRepository, txManager db.TxManager) service.UserService {
	return &serv{
		userRepo:  userRepo,
		txManager: txManager,
	}
}
