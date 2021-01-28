package routes

import (
	"github.com/gofiber/fiber/v2"
)

func WebRoutes(web fiber.Router) {
	LandingRoutes(web)
	WebAuthRoutes(web)
	UserRoutes(web)
}
