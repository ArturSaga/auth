package service_error

import "errors"

var (
	ErrPasswordsNotMatch = errors.New("bad request data, validate error")
	ErrCreateUser        = errors.New("failed to create user")
	ErrGetUser           = errors.New("failed to get user")
	ErrUpdateUser        = errors.New("failed to update user")
	ErrGrpcHostNotFound  = errors.New("grpc host not found")
	ErrPgDsnNotFound     = errors.New("pg dsn not found")
)
