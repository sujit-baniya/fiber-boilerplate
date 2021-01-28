package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"
	"github.com/sujit-baniya/verify-rest/app"
	"github.com/sujit-baniya/verify-rest/pkg/auth"
	"github.com/sujit-baniya/verify-rest/pkg/models"
)

func RedirectToHomePageOnLogin(c *fiber.Ctx) error {
	if auth.IsLoggedIn(c) {
		return c.Redirect("/")
	}
	return c.Next()
}

func ValidateLoginPost(c *fiber.Ctx) error {
	var login models.Login
	if err := c.BodyParser(&login); err != nil {
		return app.Http.Flash.WithError(c, fiber.Map{
			"message": err.Error(),
		}).Redirect("/login")
	}
	v := validate.Struct(login)
	if !v.Validate() {
		return app.Http.Flash.WithError(c, fiber.Map{
			"message": v.Errors.One(),
		}).Redirect("/login")
	}
	user, err := login.CheckLogin() //nolint:wsl

	if err != nil {
		return app.Http.Flash.WithError(c, fiber.Map{
			"message": err.Error(),
		}).Redirect("/login")
	}
	c.Locals("user", user)
	return c.Next()
}
