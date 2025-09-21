package temotes

import (
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

func getEnvOrFile(key string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	if filePath := os.Getenv(key + "_FILE"); filePath != "" {
		data, err := os.ReadFile(filePath)
		if err == nil {
			return strings.TrimSpace(string(data))
		}
	}
	return ""
}

func getInt(key string, def int) int {
	val := getEnvOrFile(key)
	if val == "" {
		return def
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return def
	}
	return i
}

type Config struct {
	ServerAddr            string
	TwitchClientID        string
	TwitchClientSecret    string
	DashboardLogin        string
	DashboardPassword     string
	ProxyMaxAge           int
	RateLimitMax          int
	RateLimitTime         int
	RedisAddr             string
	RedisPassword         string
	RedisDB               int
	CacheTTLGlobalEmotes  int
	CacheTTLChannelEmotes int
	CacheTTLTwitchIDs     int
	FetcherTimeout        int
}

func LoadConfig() {
	cfg := &Config{
		ServerAddr:            getEnvOrFile("SERVER_ADDR"),
		TwitchClientID:        getEnvOrFile("TWITCH_CLIENT_ID"),
		TwitchClientSecret:    getEnvOrFile("TWITCH_CLIENT_SECRET"),
		DashboardLogin:        getEnvOrFile("DASHBOARD_LOGIN"),
		DashboardPassword:     getEnvOrFile("DASHBOARD_PASSWORD"),
		ProxyMaxAge:           getInt("PROXY_MAX_AGE", 86400),
		RateLimitMax:          getInt("RATE_LIMIT_MAX", 10),
		RateLimitTime:         getInt("RATE_LIMIT_TIME", 10),
		RedisAddr:             getEnvOrFile("REDIS_ADDR"),
		RedisPassword:         getEnvOrFile("REDIS_PASSWORD"),
		RedisDB:               getInt("REDIS_DB", 0),
		CacheTTLGlobalEmotes:  getInt("CACHE_TTL_GLOBAL_EMOTES", 900),
		CacheTTLChannelEmotes: getInt("CACHE_TTL_CHANNEL_EMOTES", 300),
		CacheTTLTwitchIDs:     getInt("CACHE_TTL_TWITCH_IDENTIFIERS", 604800),
		FetcherTimeout:        getInt("FETCHER_TIMEOUT", 3),
	}
	SetConfig(cfg)
	log.Print("Config has been loaded.")
}

var GlobalConfig *Config

func SetConfig(cfg *Config) {
	GlobalConfig = cfg
}

var once sync.Once

func GetConfig() *Config {
	once.Do(LoadConfig)
	return GlobalConfig
}
