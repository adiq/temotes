package endpoints

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"strings"
	"temotes/temotes"
	"temotes/temotes/providers"
	"time"
)

type Helpers struct{}

func (h Helpers) GetTwitchUserId(input string) (temotes.TwitchUserId, error) {
	id, err := strconv.ParseInt(strings.ToLower(input), 10, 64)

	if id == 0 || err != nil {
		cache := temotes.CacheService{}.GetCacheClient()
		cacheKey := fmt.Sprintf("twitch_channel_id_%s", strings.ToLower(input))
		channelId, channelIdCached := cache.Get(cacheKey)

		if channelIdCached {
			return channelId.(temotes.TwitchUserId), nil
		}

		twitchUser, err := providers.TwitchFetcher{}.FetchUserIdentifiers(input)
		if err != nil {
			return 0, fiber.NewError(fiber.StatusNotFound, "User not found")
		}

		cache.Set(cacheKey, twitchUser.ID, time.Hour)

		return twitchUser.ID, nil
	}

	return temotes.TwitchUserId(id), nil
}

func (h Helpers) ParseServices(input string) (*[]temotes.EmoteFetcher, error) {
	var fetchers []temotes.EmoteFetcher

	unfilteredServices := strings.Split(strings.ToLower(input), ".")
	services := temotes.Unique(unfilteredServices)

	if len(services) == 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "No services specified")
	}

	if len(services) != len(unfilteredServices) {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Duplicate services specified")
	}

	if temotes.Contains(services, "all") {
		if len(services) > 1 {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Cannot specify all and other services")
		}

		services = []string{"twitch", "7tv", "bttv", "ffz"}
	}

	for _, serviceName := range services {
		switch serviceName {
		case "twitch":
			fetchers = append(fetchers, providers.TwitchFetcher{})
		case "7tv":
			fetchers = append(fetchers, providers.SevenTvFetcher{})
		case "bttv":
			fetchers = append(fetchers, providers.BttvFetcher{})
		case "ffz":
			fetchers = append(fetchers, providers.FfzFetcher{})
		default:
			return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Invalid service: %s", serviceName))
		}
	}

	return &fetchers, nil
}
