package model

import (
	"time"

	desc "github.com/ArturSaga/auth/api/grpc/pkg/user_v1"
)

// User - модель, для работы со слоем репозитория
type User struct {
	ID        int64     `db:"id"`
	Info      UserInfo  `db:""`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// UserInfo - модель, для работы со слоем репозитория
type UserInfo struct {
	Name            string    `db:"name"`
	Email           string    `db:"email"`
	Password        string    `db:"password_hash"`
	PasswordConfirm string    `db:""`
	Role            desc.Role `db:"role"`
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
