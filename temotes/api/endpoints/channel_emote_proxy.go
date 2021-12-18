package endpoints

import (
	"github.com/gofiber/fiber/v2"
	"temotes/temotes"
)

func GetChannelEmoteProxy(c *fiber.Ctx) error {
	emoteCode := c.Params("emote")
	if emoteCode == "" {
		return fiber.NewError(fiber.StatusBadRequest, "emote is required")
	}

	rawSize := c.Params("size")
	if rawSize == "" {
		return fiber.NewError(fiber.StatusBadRequest, "size is required")
	}

	var size temotes.EmoteSize
	for _, emSize := range []temotes.EmoteSize{temotes.Size1x, temotes.Size2x, temotes.Size3x, temotes.Size4x} {
		if emSize == temotes.EmoteSize(rawSize) {
			size = emSize
		}
	}

	if &size == nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid size")
	}

	channelId, err := Helpers{}.GetTwitchUserId(c.Params("channel"))
	if err != nil {
		return err
	}

	fetchers, err := Helpers{}.ParseServices(c.Params("services"))
	if err != nil {
		return err
	}

	for _, fetcher := range *fetchers {
		emotes := fetcher.FetchChannelEmotes(channelId)
		for _, emote := range emotes {
			if emote.Code == emoteCode {
				for _, url := range emote.Urls {
					if url.Size == size {
						return c.Redirect(url.Url, fiber.StatusTemporaryRedirect)
					}
				}
			}
		}
	}

	return fiber.NewError(fiber.StatusNotFound, "emote not found")
}
