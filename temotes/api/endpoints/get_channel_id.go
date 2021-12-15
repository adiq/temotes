package endpoints

import (
	"github.com/gofiber/fiber/v2"
	"strings"
	"temotes/temotes/services"
)

func GetChannelId(c *fiber.Ctx) error {
	channel := strings.ToLower(c.Params("channel"))
	channelId, err := services.GetTwitchUserId(channel)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	return c.JSON(fiber.Map{
		"login": channel,
		"id":    channelId,
	})
}
