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

	implementation *user.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

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

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewUserRepository(s.DBClient(ctx))
	}

	return s.userRepository
}

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userServ == nil {
		s.userServ = userService.NewUserService(
			s.UserRepository(ctx),
			s.TxManager(ctx),
		)
	}

	return s.userServ
}

func (s *serviceProvider) UserImpl(ctx context.Context) *user.Implementation {
	if s.implementation == nil {
		s.implementation = user.NewImplementation(s.UserService(ctx))
	}

	return s.implementation
}
