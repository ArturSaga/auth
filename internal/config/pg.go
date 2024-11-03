package config

import (
	"os"

	serviceErr "github.com/ArturSaga/auth/internal/service_error"
)

const (
	dsnEnvName = "PG_DSN"
)

// PGConfig - интерфейс, определяющий методы PGConfig
type PGConfig interface {
	DSN() string
}

type pgConfig struct {
	dsn string
}

// NewPGConfig - публичный метод, создающий новое подключение к Postgres
func NewPGConfig() (PGConfig, error) {
	dsn := os.Getenv(dsnEnvName)
	if len(dsn) == 0 {
		return nil, serviceErr.ErrPgDsnNotFound
	}

	return &pgConfig{
		dsn: dsn,
	}, nil
}

// DSN - публичный метод, возвращающий DSN
func (cfg *pgConfig) DSN() string {
	return cfg.dsn
}
