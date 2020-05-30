package routes

import (
	"fmt"
	"github.com/gofiber/fiber"
	Controller "github.com/thomasvvugt/fiber-boilerplate/app/controllers/web"
)

func RegisterWeb(app *fiber.App) {
	errorPages(app)
	// Homepage
	app.Get("/", Controller.Index)

	// Auth routes
	app.Get("/login", Controller.ShowLoginForm)
	app.Post("/login", Controller.PostLoginForm)
	app.Post("/logout", Controller.PostLogoutForm)
}
func errorPages(app *fiber.App) {
	// Homepage
	app.Get("/404", func(c *fiber.Ctx) {
		errorPageHandlers(c, 404)
	})
	// Homepage
	app.Get("/500", func(c *fiber.Ctx) {
		errorPageHandlers(c, 500)
	})
}

func errorPageHandlers(c *fiber.Ctx, code int) {
	c.SendStatus(code)
	c.Render(fmt.Sprintf("errors/%d", code), fiber.Map{})
}
