package tests

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	"github.com/ArturSaga/platform_common/pkg/db"

	desc "github.com/ArturSaga/auth/api/grpc/pkg/user_v1"
	"github.com/ArturSaga/auth/internal/client/cache"
	redisMocks "github.com/ArturSaga/auth/internal/client/cache/mocks"
	txMocks "github.com/ArturSaga/auth/internal/client/db/mocks"
	serviceConverter "github.com/ArturSaga/auth/internal/convertor"
	serviceModel "github.com/ArturSaga/auth/internal/model"
	"github.com/ArturSaga/auth/internal/repository"
	repoMocks "github.com/ArturSaga/auth/internal/repository/mocks"
	"github.com/ArturSaga/auth/internal/service/user"
)

func Test_serv_GetUser(t *testing.T) {
	type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository
	type redisClientMockFunc func(mc *minimock.Controller) cache.RedisClient
	type transactionMockFunc func(mc *minimock.Controller) db.TxManager

	t.Parallel()

	// Инициализация контекста и вспомогательных данных
	ctx := context.Background()
	mc := minimock.NewController(t)

	// Генерация тестовых данных с помощью gofakeit
	id := int64(1)
	name := gofakeit.Name()
	email := gofakeit.Email()
	password := gofakeit.Word()
	time := gofakeit.Date()
	role := desc.Role_ADMIN
	//roleString := desc.Role_ADMIN.String()
	serviceErr := errors.New("service error")
	emptyValues := make([]interface{}, 0)

	//redisUserModel := &redisModel.User{
	//	ID:              id,
	//	Name:            name,
	//	Email:           email,
	//	Password:        password,
	//	PasswordConfirm: "",
	//	Role:            roleString,
	//	CreatedAt:       1731567461,
	//	UpdatedAt:       1731567461,
	//}

	//redisInterface := []interface{}{
	//	strconv.FormatInt(redisUserModel.ID, 10), // Преобразуем ID в строку
	//	redisUserModel.Name,
	//	redisUserModel.Email,
	//	redisUserModel.Password,
	//	redisUserModel.PasswordConfirm,
	//	redisUserModel.Role,
	//	strconv.FormatInt(redisUserModel.CreatedAt, 10), // Преобразуем CreatedAt в строку
	//	strconv.FormatInt(redisUserModel.UpdatedAt, 10), // Преобразуем UpdatedAt в строку
	//}

	modelUserInfo := serviceModel.UserInfo{
		Name:            name,
		Email:           email,
		Password:        password,
		PasswordConfirm: "",
		Role:            role,
	}

	servUserModel := &serviceModel.User{
		ID:        id,
		Info:      modelUserInfo,
		CreatedAt: time,
		UpdatedAt: time,
	}

	// Очистка контроллера после завершения теста
	t.Cleanup(mc.Finish)

	// Группировка тестов
	tests := []struct {
		name string
		args struct {
			ctx context.Context
			id  int64
		}
		want           *serviceModel.User
		err            error
		userRepository userRepositoryMockFunc
		redisClient    redisClientMockFunc
		txManager      transactionMockFunc
	}{
		{
			name: "success get from repo",
			args: struct {
				ctx context.Context
				id  int64
			}{
				ctx: ctx,
				id:  id,
			},
			want: servUserModel,
			err:  nil,
			redisClient: func(mc *minimock.Controller) cache.RedisClient {
				mock := redisMocks.NewRedisClientMock(mc)
				mock.HGetAllMock.Expect(ctx, strconv.FormatInt(id, 10)).Return(emptyValues, nil)
				mock.HashSetMock.Expect(ctx, strconv.FormatInt(id, 10), serviceConverter.ToUserRedisFromService(servUserModel)).Return(nil)
				return mock
			},
			userRepository: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.GetUserMock.Expect(ctx, id).Return(servUserModel, nil)
				return mock
			},
			txManager: func(mc *minimock.Controller) db.TxManager {
				mock := txMocks.NewTxManagerMock(mc)
				return mock
			},
		},
		{
			name: "failed get from repo",
			args: struct {
				ctx context.Context
				id  int64
			}{
				ctx: ctx,
				id:  id,
			},
			want: nil,
			err:  serviceErr,
			redisClient: func(mc *minimock.Controller) cache.RedisClient {
				mock := redisMocks.NewRedisClientMock(mc)
				mock.HGetAllMock.Expect(ctx, strconv.FormatInt(id, 10)).Return(emptyValues, nil)
				return mock
			},
			userRepository: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.GetUserMock.Expect(ctx, id).Return(nil, serviceErr)
				return mock
			},
			txManager: func(mc *minimock.Controller) db.TxManager {
				mock := txMocks.NewTxManagerMock(mc)
				return mock
			},
		},
		{
			name: "failed get from cache",
			args: struct {
				ctx context.Context
				id  int64
			}{
				ctx: ctx,
				id:  id,
			},
			want: nil,
			err:  serviceErr,
			redisClient: func(mc *minimock.Controller) cache.RedisClient {
				mock := redisMocks.NewRedisClientMock(mc)
				mock.HGetAllMock.Expect(ctx, strconv.FormatInt(id, 10)).Return(nil, serviceErr)
				return mock
			},
			userRepository: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				return mock
			},
			txManager: func(mc *minimock.Controller) db.TxManager {
				mock := txMocks.NewTxManagerMock(mc)
				return mock
			},
		},
		{
			name: "failed set to cache",
			args: struct {
				ctx context.Context
				id  int64
			}{
				ctx: ctx,
				id:  id,
			},
			want: nil,
			err:  serviceErr,
			redisClient: func(mc *minimock.Controller) cache.RedisClient {
				mock := redisMocks.NewRedisClientMock(mc)
				mock.HGetAllMock.Expect(ctx, strconv.FormatInt(id, 10)).Return(emptyValues, nil)
				mock.HashSetMock.Expect(ctx, strconv.FormatInt(id, 10), serviceConverter.ToUserRedisFromService(servUserModel)).Return(serviceErr)
				return mock
			},
			userRepository: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.GetUserMock.Expect(ctx, id).Return(servUserModel, nil)
				return mock
			},
			txManager: func(mc *minimock.Controller) db.TxManager {
				mock := txMocks.NewTxManagerMock(mc)
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
			newID, err := userService.GetUser(tt.args.ctx, tt.args.id)

			// Проверка результата
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.want, newID)
		})
	}
}
