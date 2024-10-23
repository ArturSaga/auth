package config

import (
	"github.com/joho/godotenv"
)

// Load - публичный метод, который загружает данные из переданного env файла
func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}

// GRPCConfig - интерфейс, определящий методы для формирования адресса работы с GRPC
type GRPCConfig interface {
	Address() string
}

// PGConfig - интерфейс, определящий методы для формирования DSN строки, для подключения к бд
type PGConfig interface {
	DSN() string
}
