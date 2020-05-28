package providers

import "github.com/gofiber/fiber"

// Force HTTPS protocol if not forwarded using a reverse proxy
func ForceHTTPS(c *fiber.Ctx)  {
	if c.Get("X-Forwarded-Proto") != "https" && c.Protocol() == "http" {
		c.Redirect("https://" + c.Hostname() + c.OriginalURL(), 308)
	}
}
