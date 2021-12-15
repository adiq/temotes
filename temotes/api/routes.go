package api

import (
	"github.com/gofiber/fiber/v2"
	"temotes/temotes/api/endpoints"
)

func setupRoutes(app *fiber.App) {
	v1 := app.Group("/v1")

	v1.Get("/channel/:channel/id", endpoints.GetChannelId)
}
