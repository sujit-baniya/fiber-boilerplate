package config

import "github.com/gofiber/fiber/v2"

func LoadHeaders(c *fiber.Ctx) error {
	// Set some security headers:
	c.Set("X-XSS-Protection", "1; mode=block")
	c.Set("X-Content-Type-Options", "nosniff")
	c.Set("X-Download-Options", "noopen")
	c.Set("Strict-Transport-Security", "max-age=5184000")
	c.Set("X-Frame-Options", "SAMEORIGIN")
	return c.Next()
}

func LoadCacheHeaders(c *fiber.Ctx) error {
	c.Set("X-DNS-Prefetch-Control", "off")
	c.Set("Pragma", "no-cache")
	c.Set("Expires", "Fri, 01 Jan 1990 00:00:00 GMT")
	c.Set("Cache-Control", "no-cache, must-revalidate, no-store, max-age=0, private")
	return c.Next()
}
