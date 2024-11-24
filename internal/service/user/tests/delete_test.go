package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ArturSaga/platform_common/pkg/db"

	"github.com/ArturSaga/auth/internal/client/cache"
	redisMocks "github.com/ArturSaga/auth/internal/client/cache/mocks"
	txMocks "github.com/ArturSaga/auth/internal/client/db/mocks"
	"github.com/ArturSaga/auth/internal/repository"
	repoMocks "github.com/ArturSaga/auth/internal/repository/mocks"
	"github.com/ArturSaga/auth/internal/service/user"
)

func Test_serv_DeleteUser(t *testing.T) {
	type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository
	type redisClientMockFunc func(mc *minimock.Controller) cache.RedisClient
	type transactionMockFunc func(mc *minimock.Controller) db.TxManager

	t.Parallel()

	// Инициализация контекста и вспомогательных данных
	ctx := context.Background()
	mc := minimock.NewController(t)

	// Генерация тестовых данных с помощью gofakeit
	id := gofakeit.Int64()
	serviceErr := fmt.Errorf("service error")

	// Очистка контроллера после завершения теста
	t.Cleanup(mc.Finish)

	// Группировка тестов
	tests := []struct {
		name string
		args struct {
			ctx context.Context
			id  int64
		}
		want           emptypb.Empty
		err            error
		userRepository userRepositoryMockFunc
		redisClient    redisClientMockFunc
		txManager      transactionMockFunc
	}{
		{
			name: "success delete",
			args: struct {
				ctx context.Context
				id  int64
			}{
				ctx: ctx,
				id:  id,
			},
			want: emptypb.Empty{},
			err:  nil,
			userRepository: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.DeleteUserMock.Expect(ctx, id).Return(emptypb.Empty{}, nil)
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
			name: "failed delete",
			args: struct {
				ctx context.Context
				id  int64
			}{
				ctx: ctx,
				id:  id,
			},
			want: emptypb.Empty{},
			err:  serviceErr,
			userRepository: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.DeleteUserMock.Expect(ctx, id).Return(emptypb.Empty{}, serviceErr)
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
			newID, err := userService.DeleteUser(tt.args.ctx, tt.args.id)

			// Проверка результата
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.want, newID)
		})
	}
}
