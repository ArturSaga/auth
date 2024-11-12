package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	desc "github.com/ArturSaga/auth/api/grpc/pkg/user_v1"
	"github.com/ArturSaga/auth/internal/api/user"
	"github.com/ArturSaga/auth/internal/model"
	"github.com/ArturSaga/auth/internal/service"
	serviceMocks "github.com/ArturSaga/auth/internal/service/mocks"
	"github.com/ArturSaga/auth/internal/service_error"
)

func TestUserApi_CreateUser(t *testing.T) {
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
	name := gofakeit.Name()
	email := gofakeit.Email()
	password := gofakeit.Word()
	passwordConfirm := gofakeit.Word()
	serviceErr := fmt.Errorf("service error")

	// Успешный случай (пароли совпадают)
	userInfoSuccess := &desc.UserInfo{
		Name:            name,
		Email:           email,
		Password:        password,
		PasswordConfirm: password,
		Role:            desc.Role_ADMIN,
	}

	// Случай с ошибкой (пароли не совпадают)
	userInfoFailedCase1 := &desc.UserInfo{
		Name:            name,
		Email:           email,
		Password:        password,
		PasswordConfirm: passwordConfirm,
		Role:            desc.Role_ADMIN,
	}

	//пустое имя
	userInfoFailedCase2 := &desc.UserInfo{
		Name:            "",
		Email:           email,
		Password:        password,
		PasswordConfirm: password,
		Role:            desc.Role_ADMIN,
	}

	//пустой email
	userInfoFailedCase3 := &desc.UserInfo{
		Name:            name,
		Email:           "",
		Password:        password,
		PasswordConfirm: password,
		Role:            desc.Role_ADMIN,
	}

	//пустой пароль
	userInfoFailedCase4 := &desc.UserInfo{
		Name:            name,
		Email:           email,
		Password:        "",
		PasswordConfirm: password,
		Role:            desc.Role_ADMIN,
	}

	reqSuccess := &desc.CreateUserRequest{
		Info: userInfoSuccess,
	}

	reqFailedCase1 := &desc.CreateUserRequest{
		Info: userInfoFailedCase1,
	}

	reqFailedCase2 := &desc.CreateUserRequest{
		Info: userInfoFailedCase2,
	}

	reqFailedCase3 := &desc.CreateUserRequest{
		Info: userInfoFailedCase3,
	}

	reqFailedCase4 := &desc.CreateUserRequest{
		Info: userInfoFailedCase4,
	}

	// Модели для сервисов
	userModelSuccess := &model.UserInfo{
		Name:            name,
		Email:           email,
		Password:        password,
		PasswordConfirm: password,
		Role:            desc.Role_ADMIN,
	}

	//userModelFailedCase1 := &model.UserInfo{
	//	Name:            name,
	//	Email:           email,
	//	Password:        password,
	//	PasswordConfirm: passwordConfirm,
	//	Role:            desc.Role_ADMIN,
	//}

	res := &desc.CreateUserResponse{
		Id: id,
	}

	// Очистка контроллера после завершения теста
	t.Cleanup(mc.Finish)

	// Группировка тестов
	tests := []struct {
		name string
		args struct {
			ctx context.Context
			req *desc.CreateUserRequest
		}
		want            *desc.CreateUserResponse
		err             error
		userServiceMock func(mc *minimock.Controller) service.UserService
	}{
		{
			name: "success create",
			args: struct {
				ctx context.Context
				req *desc.CreateUserRequest
			}{
				ctx: ctx,
				req: reqSuccess,
			},
			want: res,
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.CreateUserMock.Expect(ctx, userModelSuccess).Return(id, nil)
				return mock
			},
		},
		{
			name: "failed create case 1",
			args: struct {
				ctx context.Context
				req *desc.CreateUserRequest
			}{
				ctx: ctx,
				req: reqFailedCase1,
			},
			want: nil,
			err:  service_error.ErrPasswordsNotMatch,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				return mock
			},
		},
		{
			name: "failed create case 2",
			args: struct {
				ctx context.Context
				req *desc.CreateUserRequest
			}{
				ctx: ctx,
				req: reqFailedCase2,
			},
			want: nil,
			err:  service_error.ErrRequireParam,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				return mock
			},
		},
		{
			name: "failed create case 3",
			args: struct {
				ctx context.Context
				req *desc.CreateUserRequest
			}{
				ctx: ctx,
				req: reqFailedCase3,
			},
			want: nil,
			err:  service_error.ErrRequireParam,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				return mock
			},
		},
		{
			name: "failed create case 4",
			args: struct {
				ctx context.Context
				req *desc.CreateUserRequest
			}{
				ctx: ctx,
				req: reqFailedCase4,
			},
			want: nil,
			err:  service_error.ErrRequireParam,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				return mock
			},
		},
		{
			name: "failed create case 5",
			args: struct {
				ctx context.Context
				req *desc.CreateUserRequest
			}{
				ctx: ctx,
				req: reqSuccess,
			},
			want: nil,
			err:  serviceErr,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.CreateUserMock.Expect(ctx, userModelSuccess).Return(0, serviceErr)
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
			newID, err := api.CreateUser(tt.args.ctx, tt.args.req)

			// Проверка результата
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.want, newID)
		})
	}
}
