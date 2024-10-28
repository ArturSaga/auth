package config

import (
	"github.com/joho/godotenv"
)

// Load - публичный метод, загружающий параметры из env файла
func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}
