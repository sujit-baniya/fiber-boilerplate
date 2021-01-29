package routes

import (
	"github.com/gofiber/fiber/v2"
	apiControllers "github.com/sujit-baniya/fiber-boilerplate/rest/controllers/api"
)

func ApiRoutes(api fiber.Router) {
	v1Routes(api)
}

func v1AuthRoutes(api fiber.Router) {
	api.Post("/oauth/token", apiControllers.OAuthToken)
}

func v1Routes(api fiber.Router) {
	v1 := api.Group("v1")
	v1AuthRoutes(v1)
}
