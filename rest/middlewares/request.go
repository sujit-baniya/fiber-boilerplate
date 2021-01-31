package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"time"
)

func MaxBodySize(sizeInMB int) fiber.Handler {
	sizeInMB = sizeInMB * 1024 * 1024
	return func(c *fiber.Ctx) error {
		if len(c.Body()) >= sizeInMB {
			// custom response here
			return fiber.ErrRequestEntityTooLarge
		}
		return c.Next()
	}
}

func Limit(maxRequest int, duration time.Duration) func(*fiber.Ctx) error {
	return limiter.New(limiter.Config{
		Max:        maxRequest,
		Expiration: duration * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(429).JSON(fiber.Map{
				"error":   true,
				"message": "Too many requests",
			})
		},
	})
}
