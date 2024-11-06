package config

import (
	"net"
	"os"

	serviceErr "github.com/ArturSaga/auth/internal/service_error"
)

const (
	grpcHostEnvName = "GRPC_HOST"
	grpcPortEnvName = "GRPC_PORT"
)

// GRPCConfig - интерфейс, определяющий методы GRPCConfig
type GRPCConfig interface {
	Address() string
}

type grpcConfig struct {
	host string
	port string
}

// NewGRPCConfig - публичный метод, для создания grpc сервера
func NewGRPCConfig() (GRPCConfig, error) {
	host := os.Getenv(grpcHostEnvName)
	if len(host) == 0 {
		return nil, serviceErr.ErrGrpcHostNotFound
	}

	port := os.Getenv(grpcPortEnvName)
	if len(port) == 0 {
		return nil, serviceErr.ErrGrpcHostNotFound
	}

	return &grpcConfig{
		host: host,
		port: port,
	}, nil
}

// Address -  публичный метод, формирующий url + port подключения к бд
func (cfg *grpcConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
