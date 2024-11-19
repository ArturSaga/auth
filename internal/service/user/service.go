package user

import (
	"github.com/ArturSaga/platform_common/pkg/db"

	"github.com/ArturSaga/auth/internal/client/cache"
	"github.com/ArturSaga/auth/internal/repository"
	"github.com/ArturSaga/auth/internal/service"
)

type serv struct {
	userRepo  repository.UserRepository
	cache     cache.RedisClient
	txManager db.TxManager
}

// NewUserService - публчиный метод, создающий сущность, для работы с сервисным слоем
func NewUserService(
	userRepo repository.UserRepository,
	cache cache.RedisClient,
	txManager db.TxManager,
) service.UserService {
	return &serv{
		userRepo:  userRepo,
		cache:     cache,
		txManager: txManager,
	}
}
