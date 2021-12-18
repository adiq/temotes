package endpoints

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"temotes/temotes"
)

func GetChannelEmoteProxy(c *fiber.Ctx) error {
	emoteCode := c.Query("emote")
	if emoteCode == "" {
		return fiber.NewError(fiber.StatusBadRequest, "emote is required")
	}

	rawSize := c.Query("size")
	var size temotes.EmoteSize
	if rawSize != "" {
		for _, emSize := range []temotes.EmoteSize{temotes.Size4x, temotes.Size3x, temotes.Size2x, temotes.Size1x} {
			if emSize == temotes.EmoteSize(rawSize) {
				size = emSize
			}
		}

		if &size == nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid size")
		}
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
		emote, err := getEmoteWithCode(emotes, emoteCode)
		if err != nil {
			continue
		}

		var url string
		if &size != nil {
			url, err = getEmoteUrlForSize(emote, size)
			if err != nil {
				continue
			}
		} else {
			url, err = getHighestAvailableEmoteSizeUrl(emote)
			if err != nil {
				continue
			}
		}

		if url != "" {
			return c.Redirect(url, fiber.StatusTemporaryRedirect)
		}
	}

	return fiber.NewError(fiber.StatusNotFound, "emote not found")
}

func getEmoteWithCode(emotes []temotes.Emote, code string) (temotes.Emote, error) {
	for _, emote := range emotes {
		if emote.Code == code {
			return emote, nil
		}
	}

	return temotes.Emote{}, errors.New("emote not found")
}

func getEmoteUrlForSize(emote temotes.Emote, size temotes.EmoteSize) (string, error) {
	for _, url := range emote.Urls {
		if url.Size == size {
			return url.Url, nil
		}
	}

	return "", errors.New("emote in specified size not found")
}

func getHighestAvailableEmoteSizeUrl(emote temotes.Emote) (string, error) {
	for _, size := range []temotes.EmoteSize{temotes.Size4x, temotes.Size3x, temotes.Size2x, temotes.Size1x} {
		url, err := getEmoteUrlForSize(emote, size)
		if err == nil {
			return url, nil
		}
	}

	return "", errors.New("emote seems to be unavailable in any size")
}
