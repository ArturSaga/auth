package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/emptypb"

	desc "github.com/ArturSaga/auth/api/grpc/pkg/user_v1"
	"github.com/ArturSaga/auth/internal/api/user"
	"github.com/ArturSaga/auth/internal/service"
	serviceMocks "github.com/ArturSaga/auth/internal/service/mocks"
)

func TestUserApi_DeleteUser(t *testing.T) {
	type fields struct {
		UnimplementedUserApiServer desc.UnimplementedUserApiServer
		userService                service.UserService
	}
	type args struct {
		ctx context.Context
		req *desc.CreateUserRequest
	}

	t.Parallel()

	// Инициализация контекста и вспомогательных данных
	ctx := context.Background()
	mc := minimock.NewController(t)

	// Генерация тестовых данных с помощью gofakeit
	id := gofakeit.Int64()
	serviceErr := fmt.Errorf("service error")
	req := &desc.DeleteUserRequest{
		Id: id,
	}

	// Очистка контроллера после завершения теста
	t.Cleanup(mc.Finish)

	// Группировка тестов
	tests := []struct {
		name string
		args struct {
			ctx context.Context
			req *desc.DeleteUserRequest
		}
		want            *emptypb.Empty
		err             error
		userServiceMock func(mc *minimock.Controller) service.UserService
	}{
		{
			name: "success delete",
			args: struct {
				ctx context.Context
				req *desc.DeleteUserRequest
			}{
				ctx: ctx,
				req: req,
			},
			want: &emptypb.Empty{},
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.DeleteUserMock.Expect(ctx, req.Id).Return(emptypb.Empty{}, nil)
				return mock
			},
		},
		{
			name: "failed delete",
			args: struct {
				ctx context.Context
				req *desc.DeleteUserRequest
			}{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				// Передаем те же данные, но с другими параметрами для failed case
				mock.DeleteUserMock.Expect(ctx, req.Id).Return(emptypb.Empty{}, serviceErr)
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
			userServiceMock := tt.userServiceMock(mc)

			// Создаем объект API
			api := user.NewUserAPI(userServiceMock)

			// Вызов функции CreateUser
			newID, err := api.DeleteUser(tt.args.ctx, tt.args.req)

			// Проверка результата
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.want, newID)
		})
	}
}
