package app

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	desc "github.com/ArturSaga/auth/api/grpc/pkg/user_v1"
	"github.com/ArturSaga/auth/internal/closer"
	"github.com/ArturSaga/auth/internal/config"
)

// App - сущность приложения, для запуска сервера и инициализации его зависимостей
type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
}

// NewApp - публичный метод, создающий сущность приложения
func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// Run - публичный метод, запускающий grpc сервер
func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	return a.runGRPCServer()
}

// initDeps - приватный метод, вызывающий инициализацию зависимостей, загрузку конфигов
func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

// initConfig - приватный метод, вызывающий загрузку конфигов из переменных окружения
func (a *App) initConfig(_ context.Context) error {
	err := config.Load("local.env")
	if err != nil {
		return err
	}

	return nil
}

// initServiceProvider - приватный метод, вызывающий инициализацию зависимостей (DI котейнер)
func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

// initGRPCServer - приватный метод, инициализирующий grpc сервер
func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

	reflection.Register(a.grpcServer)

	desc.RegisterUserApiServer(a.grpcServer, a.serviceProvider.UserImpl(ctx))

	return nil
}

// runGRPCServer - приватный метод, запускающий grpc сервер
func (a *App) runGRPCServer() error {
	log.Printf("GRPC server is running on %s", a.serviceProvider.GRPCConfig().Address())

	list, err := net.Listen("tcp", a.serviceProvider.GRPCConfig().Address())
	if err != nil {
		return err
	}

	err = a.grpcServer.Serve(list)
	if err != nil {
		return err
	}

	return nil
}
