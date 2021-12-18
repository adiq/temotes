package temotes

import "github.com/go-redis/redis/v8"

var redisClient *redis.Client

type CacheService struct{}

func (c CacheService) Connect() *redis.Client {
	if redisClient == nil {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})
	}

	return redisClient
}
