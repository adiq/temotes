package temotes

import (
	"context"
	"log"
	"net/http"
	"time"
)

const (
	GlobalEmotesTtl  = time.Second * 5
	ChannelEmotesTtl = time.Second * 5
	TwitchIdTtl      = time.Second * 5
)

var fetcher *Fetcher

type CachedFetcher struct{}

func (f CachedFetcher) getFetcher() *Fetcher {
	if fetcher == nil {
		fetcher = &Fetcher{}
	}

	return fetcher
}

func (f CachedFetcher) FetchData(url string, ttl time.Duration, cacheKey string) []byte {
	req := f.getFetcher().GetRequest(url)

	return f.FetchDataRequest(req, ttl, cacheKey)
}

func (f CachedFetcher) FetchDataRequest(req *http.Request, ttl time.Duration, cacheKey string) []byte {
	cache := CacheService{}.Connect()
	cacheData, err := cache.Get(context.Background(), cacheKey).Result()
	if err == nil {
		log.Println("Cache hit for " + cacheKey)
		return []byte(cacheData)
	}

	log.Println("Cache miss for " + cacheKey)
	body := f.getFetcher().FetchDataRequest(req)

	go cache.Set(context.Background(), cacheKey, string(body), ttl)

	return body
}
