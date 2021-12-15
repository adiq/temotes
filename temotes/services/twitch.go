package services

import (
	"strconv"
	"temotes/temotes"
	"temotes/temotes/providers"
)

func GetTwitchUserId(input string) (temotes.TwitchUserId, error) {
	id, err := strconv.ParseInt(input, 10, 64)

	if id == 0 || err != nil {
		userId, userErr := providers.TwitchFetcher{}.FetchUserId(input)
		if userErr != nil {
			return temotes.TwitchUserId(0), userErr
		}

		return userId, nil
	}

	return temotes.TwitchUserId(id), nil
}
