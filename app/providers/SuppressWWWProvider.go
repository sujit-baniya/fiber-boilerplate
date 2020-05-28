package providers

import (
	"strings"

	"github.com/gofiber/fiber"
)

// Suppress the `www.` at the beginning of URLs
func SuppressWWW(c *fiber.Ctx)  {
	hostnameSplit := strings.Split(c.Hostname(), ".")
	if hostnameSplit[0] == "www" && len(hostnameSplit) > 1 {
		newHostname := ""
		for i := 1; i <= (len(hostnameSplit) - 1); i++ {
			if i != (len(hostnameSplit) - 1) {
				newHostname = newHostname + hostnameSplit[i] + "."
			} else {
				newHostname = newHostname + hostnameSplit[i]
			}
		}
		c.Redirect(c.Protocol() + "://" + newHostname + c.OriginalURL(), 301)
	}
}
