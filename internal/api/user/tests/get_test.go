package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"

	desc "github.com/ArturSaga/auth/api/grpc/pkg/user_v1"
	"github.com/ArturSaga/auth/internal/api/user"
	"github.com/ArturSaga/auth/internal/model"
	"github.com/ArturSaga/auth/internal/service"
	serviceMocks "github.com/ArturSaga/auth/internal/service/mocks"
	serviceErr "github.com/ArturSaga/auth/internal/service_error"
)

func TestUserApi_GetUser(t *testing.T) {
	type fields struct {
		UnimplementedUserApiServer desc.UnimplementedUserApiServer
		userService                service.UserService
	}
	type args struct {
		ctx context.Context
		req *desc.GetUserRequest
	}

	t.Parallel()

	// Инициализация контекста и вспомогательных данных
	ctx := context.Background()
	mc := minimock.NewController(t)

	// Генерация тестовых данных с помощью gofakeit
	id := gofakeit.Int64()
	name := gofakeit.Name()
	email := gofakeit.Email()
	password := gofakeit.Word()
	time := gofakeit.Date()
	err := errors.New("service error")

	modelUserInfo := model.UserInfo{
		Name:            name,
		Email:           email,
		Password:        password,
		PasswordConfirm: "",
	}

	servUserModel := &model.User{
		ID:        id,
		Info:      modelUserInfo,
		CreatedAt: time,
		UpdatedAt: time,
	}

	userApi := &desc.User{
		Id: id,
		Info: &desc.UserInfo{
			Name:            name,
			Email:           email,
			Password:        password,
			PasswordConfirm: "",
		},
		CreatedAt: timestamppb.New(time),
		UpdatedAt: timestamppb.New(time),
	}

	//err := fmt.Errorf("service error")
	req := &desc.GetUserRequest{
		Id: id,
	}

	res := &desc.GetUserResponse{
		User: userApi,
	}

	// Очистка контроллера после завершения теста
	t.Cleanup(mc.Finish)

	// Группировка тестов
	tests := []struct {
		name string
		args struct {
			ctx context.Context
			req *desc.GetUserRequest
		}
		want            *desc.GetUserResponse
		err             error
		userServiceMock func(mc *minimock.Controller) service.UserService
	}{
		{
			name: "success get",
			args: struct {
				ctx context.Context
				req *desc.GetUserRequest
			}{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.GetUserMock.Expect(ctx, req.Id).Return(servUserModel, nil)
				return mock
			},
		},
		{
			name: "failed get",
			args: struct {
				ctx context.Context
				req *desc.GetUserRequest
			}{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  err,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.GetUserMock.Expect(ctx, req.Id).Return(nil, err)
				return mock
			},
		},
		{
			name: "success but convertor returns nil",
			args: struct {
				ctx context.Context
				req *desc.GetUserRequest
			}{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr.ErrConvertUser, // Конвертер возвращает nil
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				// Мокируем успешное получение пользователя, но конвертер возвращает nil
				mock.GetUserMock.Expect(ctx, req.Id).Return(nil, serviceErr.ErrConvertUser)
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

			// Вызов функции GetUser
			newID, err := api.GetUser(tt.args.ctx, tt.args.req)

			// Проверка результата
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.want, newID)
		})
	}
}
