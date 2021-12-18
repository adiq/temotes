package temotes

import (
	"context"
	"net/http"
	"time"
)

const (
	GlobalEmotesTtl  = time.Minute * 15
	ChannelEmotesTtl = time.Minute * 5
	TwitchIdTtl      = time.Hour * 24 * 7
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
		return []byte(cacheData)
	}

	body := f.getFetcher().FetchDataRequest(req)

	go cache.Set(context.Background(), cacheKey, string(body), ttl)

	return body
}
