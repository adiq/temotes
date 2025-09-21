package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"temotes/temotes"
)

func setupDashboard(app *fiber.App) {
	app.Get("/dashboard", basicauth.New(basicauth.Config{
		Users: map[string]string{
			temotes.GetConfig().DashboardLogin: temotes.GetConfig().DashboardPassword,
		},
	}), monitor.New())
}
