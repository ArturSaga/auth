package model

import (
	"time"

	"google.golang.org/protobuf/types/known/wrapperspb"

	desc "github.com/ArturSaga/auth/api/grpc/pkg/user_v1"
)

// User - модель, для работы с сервисным слоем
type User struct {
	ID        int64
	Info      UserInfo
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UserInfo - модель, для работы с сервисным слоем
type UserInfo struct {
	Name            string    `db:"name"`
	Email           string    `db:"email"`
	Password        string    `db:"password_hash"`
	PasswordConfirm string    `db:"password_hash"`
	Role            desc.Role `db:"role"`
}

// UpdateUserInfo - модель, для работы с сервисным слоем
type UpdateUserInfo struct {
	UserID          int64                   `db:"id"`
	Name            *wrapperspb.StringValue `db:"name"`
	OldPassword     *wrapperspb.StringValue `db:"old_password"`
	Password        *wrapperspb.StringValue `db:"password"`
	PasswordConfirm *wrapperspb.StringValue `db:"password_confirm"`
	Role            desc.Role
}

// RoleFromString Функция для преобразования Role из строкового представления
func RoleFromString(s string) desc.Role {
	switch s {
	case "UNKNOWN":
		return desc.Role_UNKNOWN
	case "ADMIN":
		return desc.Role_ADMIN
	case "USER":
		return desc.Role_USER
	default:
		return desc.Role_UNKNOWN
	}
}
