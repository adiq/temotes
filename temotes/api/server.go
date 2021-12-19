package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func SetupServer() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandler,
	})

	app.Use(cors.New())
	app.Use(recover.New())

	rateLimiter(app)
	setupRoutes(app)
	setupDashboard(app)

	app.Use(notFoundHandler)

	return app
}
