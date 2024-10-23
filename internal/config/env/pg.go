package env

import (
	"errors"
	"os"

	"github.com/ArturSaga/auth/internal/config"
)

var _ config.PGConfig = (*pgConfig)(nil)

const (
	dsnEnvName = "PG_DSN"
)

type pgConfig struct {
	dsn string
}

// NewPGConfig - публчиный метод, который формирует строку для подключения к бд, из данных env файла
func NewPGConfig() (*pgConfig, error) {
	dsn := os.Getenv(dsnEnvName)
	if len(dsn) == 0 {
		return nil, errors.New("pg dsn not found")
	}

	return &pgConfig{
		dsn: dsn,
	}, nil
}

// DSN - публчиный метод, который возвращает DSN строку для подключения к бд
func (cfg *pgConfig) DSN() string {
	return cfg.dsn
}
