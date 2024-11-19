package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	desc "github.com/ArturSaga/auth/api/grpc/pkg/user_v1"
	"github.com/ArturSaga/auth/internal/api/user"
	"github.com/ArturSaga/auth/internal/model"
	"github.com/ArturSaga/auth/internal/service"
	serviceMocks "github.com/ArturSaga/auth/internal/service/mocks"
	_ "github.com/ArturSaga/auth/internal/service_error"
	serviceErr "github.com/ArturSaga/auth/internal/service_error"
)

func TestUserApi_UpdateUser(t *testing.T) {
	t.Parallel()

	// Инициализация контекста и вспомогательных данных
	ctx := context.Background()
	mc := minimock.NewController(t)

	// Генерация тестовых данных с помощью gofakeit
	id := gofakeit.Int64()
	name := gofakeit.Name()
	email := gofakeit.Email()
	currentPassword, _ := bcrypt.GenerateFromPassword([]byte(gofakeit.Word()), 10)
	oldPassword := gofakeit.Word()
	oldPasswordHashed, _ := bcrypt.GenerateFromPassword([]byte(oldPassword), 10)
	newPassword := gofakeit.Word()
	failedPassword := gofakeit.Word()
	role := desc.Role_ADMIN
	time := time.Now()
	servErr := errors.New("service error")

	reqUpdateSuccess := &desc.UpdateUserRequest{
		Info: &desc.UpdateUserInfo{
			UserID:          id,
			Name:            wrapperspb.String(name),
			OldPassword:     wrapperspb.String(oldPassword),
			Password:        wrapperspb.String(newPassword),
			PasswordConfirm: wrapperspb.String(newPassword),
			Role:            role,
		},
	}

	reqUpdateFailed := &desc.UpdateUserRequest{
		Info: &desc.UpdateUserInfo{
			UserID:          id,
			Name:            wrapperspb.String(name),
			OldPassword:     wrapperspb.String(oldPassword),
			Password:        wrapperspb.String(newPassword),
			PasswordConfirm: wrapperspb.String(failedPassword),
			Role:            role,
		},
	}

	reqUpdateFailedCase2 := &desc.UpdateUserRequest{
		Info: &desc.UpdateUserInfo{
			UserID:          id,
			Name:            wrapperspb.String(name),
			OldPassword:     wrapperspb.String(oldPassword),
			Password:        wrapperspb.String(newPassword),
			PasswordConfirm: wrapperspb.String(newPassword),
			Role:            role,
		},
	}

	reqGet := &desc.GetUserRequest{
		Id: id,
	}

	// Модели для сервисов
	updateUserModelSuccess := &model.UpdateUserInfo{
		UserID:          id,
		Name:            &name,
		OldPassword:     &oldPassword,
		Password:        &newPassword,
		PasswordConfirm: &newPassword,
		Role:            &role,
	}

	//updateUserModelFailed := &model.UpdateUserInfo{
	//	UserID:          id,
	//	Name:            &name,
	//	OldPassword:     &oldPassword,
	//	Password:        &newPassword,
	//	PasswordConfirm: &failedPassword,
	//	Role:            &role,
	//}

	modelUserInfo := model.UserInfo{
		Name:            name,
		Email:           email,
		Password:        string(oldPasswordHashed), // Используем хэшированный старый пароль
		PasswordConfirm: "",
	}

	failedUserInfoModel := model.UserInfo{
		Name:            name,
		Email:           email,
		Password:        string(currentPassword), // Используем хэшированный старый пароль
		PasswordConfirm: "",
	}

	servUserModel := &model.User{
		ID:        id,
		Info:      modelUserInfo,
		CreatedAt: time,
		UpdatedAt: time,
	}

	failedServUserModel := &model.User{
		ID:        id,
		Info:      failedUserInfoModel,
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
			req *desc.UpdateUserRequest
		}
		want            *emptypb.Empty
		err             error
		userServiceMock func(mc *minimock.Controller) service.UserService
	}{
		{
			name: "success update",
			args: struct {
				ctx context.Context
				req *desc.UpdateUserRequest
			}{
				ctx: ctx,
				req: reqUpdateSuccess,
			},
			want: &emptypb.Empty{},
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.GetUserMock.Expect(ctx, reqGet.GetId()).Return(servUserModel, nil)
				mock.UpdateUserMock.Expect(ctx, updateUserModelSuccess).Return(emptypb.Empty{}, nil)
				return mock
			},
		},
		{
			name: "failed update",
			args: struct {
				ctx context.Context
				req *desc.UpdateUserRequest
			}{
				ctx: ctx,
				req: reqUpdateSuccess,
			},
			want: &emptypb.Empty{},
			err:  servErr,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.GetUserMock.Expect(ctx, reqGet.GetId()).Return(servUserModel, nil)
				mock.UpdateUserMock.Expect(ctx, updateUserModelSuccess).Return(emptypb.Empty{}, servErr)
				return mock
			},
		},
		{
			name: "failed update password not match",
			args: struct {
				ctx context.Context
				req *desc.UpdateUserRequest
			}{
				ctx: ctx,
				req: reqUpdateFailed,
			},
			want: &emptypb.Empty{},
			err:  serviceErr.ErrPasswordsNotMatch,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.GetUserMock.Expect(ctx, reqGet.GetId()).Return(servUserModel, nil)
				return mock
			},
		},
		{
			name: "failed update old input password not equal to current password",
			args: struct {
				ctx context.Context
				req *desc.UpdateUserRequest
			}{
				ctx: ctx,
				req: reqUpdateFailedCase2,
			},
			want: &emptypb.Empty{},
			err:  serviceErr.ErrCompareOldPassword,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.GetUserMock.Expect(ctx, reqGet.GetId()).Return(failedServUserModel, nil)
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

			// Вызов функции UpdateUser
			newID, err := api.UpdateUser(tt.args.ctx, tt.args.req)

			// Проверка результата
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.want, newID)
		})
	}
}
