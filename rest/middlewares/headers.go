package middlewares

import "github.com/gofiber/fiber/v2"

func LoadCacheHeaders(c *fiber.Ctx) error {
	c.Set("X-DNS-Prefetch-Control", "off")
	c.Set("Pragma", "no-cache")
	c.Set("Expires", "Fri, 01 Jan 1990 00:00:00 GMT")
	c.Set("Cache-Control", "no-cache, must-revalidate, no-store, max-age=0, private")
	return c.Next()
}
