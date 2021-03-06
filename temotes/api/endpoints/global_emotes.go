package endpoints

import (
	"github.com/gofiber/fiber/v2"
	"temotes/temotes"
)

func GetGlobalEmotes(c *fiber.Ctx) error {
	fetchers, err := Helpers{}.ParseServices(c.Params("services"))
	if err != nil {
		return err
	}

	var emotes []temotes.Emote
	for _, fetcher := range *fetchers {
		emotes = append(emotes, fetcher.FetchGlobalEmotes()...)
	}

	return c.JSON(emotes)
}
