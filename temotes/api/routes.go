package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"temotes/temotes/api/endpoints"
)

func setupRoutes(app *fiber.App) {
	v1 := app.Group("/v1")
	v1Global := v1.Group("/global")
	v1Channel := v1.Group("/channel/:channel")

	// Global
	v1Global.Get("/emotes/:services", endpoints.GetGlobalEmotes)

	// Channel specific
	v1Channel.Get("/emotes/:services", endpoints.GetChannelEmotes)
	v1Channel.Get("/emotes/:services/proxy", endpoints.GetChannelEmoteProxy)

	v1Channel.Get("/id", endpoints.GetChannelIdentifiers)

	// Healthcheck
	app.Use(healthcheck.New())
}
