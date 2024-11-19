package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	"github.com/ArturSaga/platform_common/pkg/db"

	desc "github.com/ArturSaga/auth/api/grpc/pkg/user_v1"
	"github.com/ArturSaga/auth/internal/client/cache"
	redisMocks "github.com/ArturSaga/auth/internal/client/cache/mocks"
	txMocks "github.com/ArturSaga/auth/internal/client/db/mocks"
	serviceModel "github.com/ArturSaga/auth/internal/model"
	"github.com/ArturSaga/auth/internal/repository"
	repoMocks "github.com/ArturSaga/auth/internal/repository/mocks"
	"github.com/ArturSaga/auth/internal/service/user"
)

func Test_serv_CreateUser(t *testing.T) {
	type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository
	type redisClientMockFunc func(mc *minimock.Controller) cache.RedisClient
	type transactionMockFunc func(mc *minimock.Controller) db.TxManager

	t.Parallel()

	// Инициализация контекста и вспомогательных данных
	ctx := context.Background()
	mc := minimock.NewController(t)

	// Генерация тестовых данных с помощью gofakeit
	id := gofakeit.Int64()
	name := gofakeit.Name()
	email := gofakeit.Email()
	password := gofakeit.Word()
	//passwordConfirm := gofakeit.Word()
	serviceErr := fmt.Errorf("service error")

	// Успешный случай (пароли совпадают)
	userInfoSuccess := &serviceModel.UserInfo{
		Name:            name,
		Email:           email,
		Password:        password,
		PasswordConfirm: password,
		Role:            desc.Role_ADMIN,
	}

	// Модели для сервисов
	userModelSuccess := &serviceModel.UserInfo{
		Name:            name,
		Email:           email,
		Password:        password,
		PasswordConfirm: password,
		Role:            desc.Role_ADMIN,
	}

	// Очистка контроллера после завершения теста
	t.Cleanup(mc.Finish)

	// Группировка тестов
	tests := []struct {
		name string
		args struct {
			ctx      context.Context
			userInfo *serviceModel.UserInfo
		}
		want           int64
		err            error
		userRepository userRepositoryMockFunc
		redisClient    redisClientMockFunc
		txManager      transactionMockFunc
	}{
		{
			name: "success create",
			args: struct {
				ctx      context.Context
				userInfo *serviceModel.UserInfo
			}{
				ctx:      ctx,
				userInfo: userInfoSuccess,
			},
			want: id,
			err:  nil,
			userRepository: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.CreateUserMock.Expect(ctx, userModelSuccess).Return(id, nil)
				return mock
			},
			txManager: func(mc *minimock.Controller) db.TxManager {
				mock := txMocks.NewTxManagerMock(mc)
				return mock
			},
			redisClient: func(mc *minimock.Controller) cache.RedisClient {
				mock := redisMocks.NewRedisClientMock(mc)
				return mock
			},
		},
		{
			name: "failed create",
			args: struct {
				ctx      context.Context
				userInfo *serviceModel.UserInfo
			}{
				ctx:      ctx,
				userInfo: userInfoSuccess,
			},
			want: 0,
			err:  serviceErr,
			userRepository: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.CreateUserMock.Expect(ctx, userModelSuccess).Return(0, serviceErr)
				return mock
			},
			txManager: func(mc *minimock.Controller) db.TxManager {
				mock := txMocks.NewTxManagerMock(mc)
				return mock
			},
			redisClient: func(mc *minimock.Controller) cache.RedisClient {
				mock := redisMocks.NewRedisClientMock(mc)
				return mock
			},
		},
	}

	// Выполнение тестов
	for _, tt := range tests {
		tt := tt // Создаем копию для каждой итерации (для параллельных тестов)
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Создаем моки для UserService
			userRepository := tt.userRepository(mc)
			redisClient := tt.redisClient(mc)
			txManager := tt.txManager(mc)

			// Создаем объект API
			userService := user.NewUserService(userRepository, redisClient, txManager)

			// Вызов функции CreateUser
			newID, err := userService.CreateUser(tt.args.ctx, tt.args.userInfo)

			// Проверка результата
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.want, newID)
		})
	}
}
