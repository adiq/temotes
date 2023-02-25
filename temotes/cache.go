package temotes

import (
	"github.com/go-redis/redis/v8"
	"github.com/patrickmn/go-cache"
	"os"
	"strconv"
	"sync"
	"time"
)

var redisClient *redis.Client
var cacheClient *cache.Cache

var redisLock = &sync.Mutex{}
var cacheLock = &sync.Mutex{}

type CacheService struct{}

func (c CacheService) GetRedisClient() *redis.Client {
	if redisClient == nil {
		redisLock.Lock()
		defer redisLock.Unlock()
		if redisClient == nil {
			redisDb, _ := strconv.ParseInt(os.Getenv("REDIS_DB"), 10, 64)

			redisClient = redis.NewClient(&redis.Options{
				Addr:     os.Getenv("REDIS_ADDR"),
				Password: os.Getenv("REDIS_PASSWORD"),
				DB:       int(redisDb),
			})
		}
	}

	return redisClient
}

func (c CacheService) GetCacheClient() *cache.Cache {
	if cacheClient == nil {
		cacheLock.Lock()
		defer cacheLock.Unlock()
		if cacheClient == nil {
			cacheClient = cache.New(time.Hour, 2*time.Hour)
		}
	}

	return cacheClient
}
