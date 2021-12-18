package endpoints

import (
	"github.com/gofiber/fiber/v2"
	"temotes/temotes/providers"
)

func GetChannelIdentifiers(c *fiber.Ctx) error {
	twitchUser, err := providers.TwitchFetcher{}.FetchUserIdentifiers(c.Params("channel"))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	return c.JSON(fiber.Map{
		"id":           twitchUser.ID,
		"login":        twitchUser.Login,
		"display_name": twitchUser.DisplayName,
	})
}
