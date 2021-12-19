package temotes

import (
	"github.com/go-redis/redis/v8"
	"os"
	"strconv"
)

var redisClient *redis.Client

type CacheService struct{}

func (c CacheService) Connect() *redis.Client {
	if redisClient == nil {
		redisDb, _ := strconv.ParseInt(os.Getenv("REDIS_DB"), 10, 64)

		redisClient = redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_ADDR"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       int(redisDb),
		})
	}

	return redisClient
}
