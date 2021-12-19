package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"os"
)

func setupDashboard(app *fiber.App) {
	app.Get("/dashboard", basicauth.New(basicauth.Config{
		Users: map[string]string{
			os.Getenv("DASHBOARD_LOGIN"): os.Getenv("DASHBOARD_PASSWORD"),
		},
	}), monitor.New())
}
