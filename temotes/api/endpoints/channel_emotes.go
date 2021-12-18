package endpoints

import (
	"github.com/gofiber/fiber/v2"
	"temotes/temotes"
)

func GetChannelEmotes(c *fiber.Ctx) error {
	channelId, err := Helpers{}.GetTwitchUserId(c.Params("channel"))
	if err != nil {
		return err
	}

	fetchers, err := Helpers{}.ParseServices(c.Params("services"))
	if err != nil {
		return err
	}

	var emotes []temotes.Emote
	for _, fetcher := range *fetchers {
		emotes = append(emotes, fetcher.FetchChannelEmotes(channelId)...)
	}

	return c.JSON(emotes)
}
