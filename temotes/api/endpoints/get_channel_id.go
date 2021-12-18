package endpoints

import (
	"github.com/gofiber/fiber/v2"
	"strings"
	"temotes/temotes/providers"
)

func GetChannelId(c *fiber.Ctx) error {
	channel := strings.ToLower(c.Params("channel"))
	twitchUser, err := providers.TwitchFetcher{}.FetchUserIdentifiers(channel)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	return c.JSON(fiber.Map{
		"id":           twitchUser.ID,
		"login":        twitchUser.Login,
		"display_name": twitchUser.DisplayName,
	})
}
