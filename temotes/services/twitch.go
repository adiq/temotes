package services

import (
	"strconv"
	"temotes/temotes"
	"temotes/temotes/providers"
)

func GetTwitchUserId(input string) temotes.TwitchUserId {
	id, err := strconv.ParseInt(input, 10, 64)

	if id == 0 || err != nil {
		return providers.TwitchFetcher{}.FetchUserId(input)
	}

	return temotes.TwitchUserId(id)
}
