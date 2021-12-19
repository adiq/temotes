package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func setupDashboard(app *fiber.App) {
	app.Get("/dashboard", basicauth.New(basicauth.Config{
		Users: map[string]string{
			"adiq": "korwin",
		},
	}), monitor.New())
}
