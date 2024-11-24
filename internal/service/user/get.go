package user

import (
	"context"
	"log"
	"strconv"

	redigo "github.com/gomodule/redigo/redis"

	redisModel "github.com/ArturSaga/auth/internal/client/cache/redis/model"
	serviceConverter "github.com/ArturSaga/auth/internal/convertor"
	serviceModel "github.com/ArturSaga/auth/internal/model"
)

// GetUser - публичный метод сервиса для получения пользователя
func (s *serv) GetUser(ctx context.Context, id int64) (*serviceModel.User, error) {
	var user redisModel.User

	// Получаем все поля из Redis
	values, err := s.cache.HGetAll(ctx, strconv.FormatInt(id, 10))
	if err != nil {
		log.Println("Error getting data from Redis:", err)
		return nil, err
	}

	// Если данных нет в кэше, получаем их из репозитория
	if len(values) == 0 {
		userRepo, err := s.userRepo.GetUser(ctx, id)
		if err != nil {
			return nil, err
		}

		userRedis := serviceConverter.ToUserRedisFromService(userRepo)
		// Кэшируем данные в Redis
		err = s.cache.HashSet(ctx, strconv.FormatInt(id, 10), userRedis)
		if err != nil {
			return nil, err
		}

		log.Printf("User data from repo: %#v", userRepo)

		return userRepo, nil
	}
	// Если данные есть в Redis, сканируем их в структуру
	err = redigo.ScanStruct(values, &user)
	if err != nil {
		log.Println("Error scanning data from Redis:", err)
		return nil, err
	}

	// Преобразуем данные из Redis в формат, используемый в сервисе
	return serviceConverter.ToUserServiceFromRedis(&user), nil
}
