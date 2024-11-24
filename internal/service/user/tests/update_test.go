package tests

import (
	"context"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/emptypb"

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

func Test_serv_UpdateUser(t *testing.T) {
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
	password := gofakeit.Word()
	role := desc.Role_ADMIN
	//serviceErr := fmt.Errorf("service error")

	updateUserInfoSuccess := &serviceModel.UpdateUserInfo{
		UserID:          id,
		Name:            &name,
		OldPassword:     &password,
		Password:        &password,
		PasswordConfirm: &password,
		Role:            &role,
	}

	//updateUserInfoFailed := &serviceModel.UpdateUserInfo{
	//	UserID:          id,
	//	Name:            &name,
	//	OldPassword:     &password,
	//	Password:        &password,
	//	PasswordConfirm: &password,
	//	Role:            &role,
	//}

	// Очистка контроллера после завершения теста
	t.Cleanup(mc.Finish)

	// Группировка тестов
	tests := []struct {
		name string
		args struct {
			ctx            context.Context
			updateUserInfo *serviceModel.UpdateUserInfo
		}
		want           emptypb.Empty
		err            error
		userRepository userRepositoryMockFunc
		redisClient    redisClientMockFunc
		txManager      transactionMockFunc
	}{
		{
			name: "success update",
			args: struct {
				ctx            context.Context
				updateUserInfo *serviceModel.UpdateUserInfo
			}{
				ctx:            ctx,
				updateUserInfo: updateUserInfoSuccess,
			},
			want: emptypb.Empty{},
			err:  nil,
			userRepository: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.UpdateUserMock.Expect(ctx, updateUserInfoSuccess).Return(emptypb.Empty{}, nil)
				return mock
			},
			txManager: func(mc *minimock.Controller) db.TxManager {
				mock := txMocks.NewTxManagerMock(mc)
				return mock
			},
			redisClient: func(mc *minimock.Controller) cache.RedisClient {
				mock := redisMocks.NewRedisClientMock(mc)
				mock.DelMock.Expect(ctx, strconv.FormatInt(updateUserInfoSuccess.UserID, 10)).Return(nil)
				return mock
			},
		},
		//{
		//	name: "failed update",
		//	args: struct {
		//		ctx            context.Context
		//		updateUserInfo *serviceModel.UpdateUserInfo
		//	}{
		//		ctx:            ctx,
		//		updateUserInfo: updateUserInfoSuccess,
		//	},
		//	want: emptypb.Empty{},
		//	err:  serviceErr,
		//	userRepository: func(mc *minimock.Controller) repository.UserRepository {
		//		mock := repoMocks.NewUserRepositoryMock(mc)
		//		mock.UpdateUserMock.Expect(ctx, updateUserInfoFailed).Return(emptypb.Empty{}, serviceErr)
		//		return mock
		//	},
		//	txManager: func(mc *minimock.Controller) db.TxManager {
		//		mock := txMocks.NewTxManagerMock(mc)
		//		return mock
		//	},
		//	redisClient: func(mc *minimock.Controller) cache.RedisClient {
		//		mock := redisMocks.NewRedisClientMock(mc)
		//		return mock
		//	},
		//},
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
			newID, err := userService.UpdateUser(tt.args.ctx, tt.args.updateUserInfo)

			// Проверка результата
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.want, newID)
		})
	}
}
