package temotes

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	fetcher          *Fetcher
	GlobalEmotesTtl  time.Duration
	ChannelEmotesTtl time.Duration
	TwitchIdTtl      time.Duration
)

type CachedFetcher struct{}

func (f CachedFetcher) initTtl() {
	globalEmotesTtl, _ := strconv.ParseInt(os.Getenv("CACHE_TTL_GLOBAL_EMOTES"), 10, 64)
	channelEmotesTtl, _ := strconv.ParseInt(os.Getenv("CACHE_TTL_CHANNEL_EMOTES"), 10, 64)
	twitchIdTtl, _ := strconv.ParseInt(os.Getenv("CACHE_TTL_TWITCH_IDENTIFIERS"), 10, 64)

	GlobalEmotesTtl = time.Duration(globalEmotesTtl) * time.Second
	ChannelEmotesTtl = time.Duration(channelEmotesTtl) * time.Second
	TwitchIdTtl = time.Duration(twitchIdTtl) * time.Second
}

func (f CachedFetcher) getFetcher() *Fetcher {
	if fetcher == nil {
		f.initTtl()
		fetcher = &Fetcher{}
	}

	return fetcher
}

func (f CachedFetcher) FetchData(url string, ttl time.Duration, cacheKey string) ([]byte, error) {
	req := f.getFetcher().GetRequest(url)

	return f.FetchDataRequest(req, ttl, cacheKey)
}

func (f CachedFetcher) FetchDataRequest(req *http.Request, ttl time.Duration, cacheKey string) ([]byte, error) {
	cache := CacheService{}.GetRedisClient()
	cacheData, err := cache.Get(context.Background(), cacheKey).Result()
	if err == nil {
		return []byte(cacheData), nil
	}

	body, err := f.getFetcher().FetchDataRequest(req)
	if err != nil {
		return nil, err
	}

	go cache.Set(context.Background(), cacheKey, string(body), ttl)

	return body, nil
}
