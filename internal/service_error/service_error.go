package service_error

import "errors"

var (
	// ErrPasswordsNotMatch - ошибка не соответствия паролей
	ErrPasswordsNotMatch = errors.New("bad request data, validate error")
	// ErrGetUser - ошибка получения данных пользователя
	ErrGetUser = errors.New("failed to get user")
	// ErrUpdateUser - ошибка при обновлении данных пользователя
	ErrUpdateUser = errors.New("failed to update user")
	// ErrGrpcHostNotFound - ошибка при получении данных grpc сервера
	ErrGrpcHostNotFound = errors.New("grpc host not found")
	// ErrPgDsnNotFound - ошибка получения данных DSN postgres
	ErrPgDsnNotFound = errors.New("pg dsn not found")
)
