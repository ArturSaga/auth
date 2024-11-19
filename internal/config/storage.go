package config

import (
	"os"

	"github.com/pkg/errors"
)

const storageModeEnvName = "STORAGE_MODE"

// StorageConfig - интрейфейс определяющий методы стораджа
type StorageConfig interface {
	Mode() string
}

// StorageConfig - структура имплементирующая интерфейс методы стораджа
type storageConfig struct {
	mode string
}

// NewStorageConfig - метод создания структура с конфигами из env
func NewStorageConfig() (*storageConfig, error) {
	storageMode := os.Getenv(storageModeEnvName)
	if len(storageMode) == 0 {
		return nil, errors.New("storage mode not found")
	}

	return &storageConfig{
		mode: storageMode,
	}, nil
}

// Mode - возврат параметра Mode
func (cfg *storageConfig) Mode() string {
	return cfg.mode
}
