package app

import (
	"context"
	"log"

	"github.com/ArturSaga/auth/internal/api/user"
	"github.com/ArturSaga/auth/internal/client/db"
	"github.com/ArturSaga/auth/internal/client/db/pg"
	"github.com/ArturSaga/auth/internal/client/db/transaction"
	"github.com/ArturSaga/auth/internal/closer"
	"github.com/ArturSaga/auth/internal/config"
	"github.com/ArturSaga/auth/internal/repository"
	userRepository "github.com/ArturSaga/auth/internal/repository/user"
	"github.com/ArturSaga/auth/internal/service"
	userService "github.com/ArturSaga/auth/internal/service/user"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	dbClient       db.Client
	txManager      db.TxManager
	userRepository repository.UserRepository

	userServ service.UserService

	userApi *user.UserAPI
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

// PgConfig - публичный метод, инициализирующий объект с postgres конфигами
func (s serviceProvider) PgConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %v", err)
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

// GRPCConfig - публичный метод, инициализирующий объект с grpc конфигами
func (s serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %v", err)
		}
		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

// DBClient - публичный метод, инициализирующий объект соединения с бд
func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PgConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

// TxManager - публичный метод, инициализирующий объект для работы с транзакциями
func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

// UserRepository - публичный метод, инициализирующий объект репозитория postgres
func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewUserRepository(s.DBClient(ctx))
	}

	return s.userRepository
}

// UserService - публичный метод, инициализирующий объект сервиса
func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userServ == nil {
		s.userServ = userService.NewUserService(
			s.UserRepository(ctx),
			s.TxManager(ctx),
		)
	}

	return s.userServ
}

// UserImpl - публичный метод, инициализирующий объект сервера
func (s *serviceProvider) UserImpl(ctx context.Context) *user.UserAPI {
	if s.userApi == nil {
		s.userApi = user.NewUserAPI(s.UserService(ctx))
	}

	return s.userApi
}
