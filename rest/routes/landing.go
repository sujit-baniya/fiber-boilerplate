package routes

import (
	"github.com/gofiber/fiber/v2"

	fibercasbinrest "github.com/prongbang/fiber-casbinrest"
	"github.com/sujit-baniya/verify-rest/app"
	"github.com/sujit-baniya/verify-rest/rest/controllers"
	"github.com/sujit-baniya/verify-rest/rest/middlewares"
)

func LandingRoutes(web fiber.Router) {
	web.Get("/", controllers.Landing)
	web.Get("/ping", fibercasbinrest.NewDefault(app.Http.Auth.Enforcer, "secret"), Pong)
	web.Get("/do/verify-email", middlewares.ValidateConfirmToken, controllers.VerifyRegisteredEmail)
}

func Pong(c *fiber.Ctx) error {
	return c.SendString("Pong")
}
