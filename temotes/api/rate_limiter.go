package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"os"
	"strconv"
	"time"
)

func rateLimiter(app *fiber.App) {
	rateLimitMax, err := strconv.ParseInt(os.Getenv("RATE_LIMIT_MAX"), 10, 64)
	if err != nil {
		rateLimitMax = 10
	}

	rateLimitTime, err := strconv.ParseInt(os.Getenv("RATE_LIMIT_TIME"), 10, 64)
	if err != nil {
		rateLimitTime = 10
	}

	app.Use(limiter.New(limiter.Config{
		KeyGenerator: func(c *fiber.Ctx) string {
			cloudflareIP := c.Get("CF-Connecting-IP")
			if cloudflareIP != "" {
				return cloudflareIP
			}

			return c.IP()
		},
		Max:        int(rateLimitMax),
		Expiration: time.Duration(rateLimitTime) * time.Second,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Rate limit reached",
			})
		},
	}))
}
