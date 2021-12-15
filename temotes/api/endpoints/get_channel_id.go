package endpoints

import (
	"github.com/gofiber/fiber/v2"
	"strings"
	"temotes/temotes/services"
)

func GetChannelId(c *fiber.Ctx) error {
	channel := strings.ToLower(c.Params("channel"))
	return c.JSON(fiber.Map{
		"login": channel,
		"id":    services.GetTwitchUserId(channel),
	})
}
