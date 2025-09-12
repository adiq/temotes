package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"temotes/temotes"
	"time"
)

func rateLimiter(app *fiber.App) {
	app.Use(limiter.New(limiter.Config{
		KeyGenerator: func(c *fiber.Ctx) string {
			cloudflareIP := c.Get("CF-Connecting-IP")
			if cloudflareIP != "" {
				return cloudflareIP
			}

			return c.IP()
		},
		Max:        temotes.GetConfig().RateLimitMax,
		Expiration: time.Duration(temotes.GetConfig().RateLimitTime) * time.Second,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Rate limit reached",
			})
		},
	}))
}
