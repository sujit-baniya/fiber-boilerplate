package routes

import (
	"github.com/gofiber/fiber"
	"log"

	Controller "github.com/thomasvvugt/fiber-boilerplate/app/controllers/web"
	"github.com/thomasvvugt/fiber-boilerplate/app/providers"
)

func RegisterWeb(app *fiber.App) {
	// Homepage
	app.Get("/", Controller.Index)

	// Panic test route, this brings up an error
	app.Get("/panic", func(c *fiber.Ctx) {
		panic("Hi, I'm a panic error!")
	})

	// Make a new hash
	app.Get("/hash/*", func(c *fiber.Ctx) {
		hash, err := providers.HashProvider().CreateHash(c.Params("*"))
		if err != nil {
			log.Fatalf("Error when creating hash: %v", err)
		}
		c.Send(hash)
	})

	// Auth routes
	app.Get("/login", Controller.ShowLoginForm)
	app.Post("/login", Controller.PostLoginForm)
	app.Post("/logout", Controller.PostLogoutForm)
}
