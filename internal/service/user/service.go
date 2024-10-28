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

func NewUserService(userRepo repository.UserRepository, txManager db.TxManager) service.UserService {
	return &serv{
		userRepo:  userRepo,
		txManager: txManager,
	}
}
