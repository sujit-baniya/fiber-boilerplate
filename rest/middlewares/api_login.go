package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"
	"github.com/sujit-baniya/verify-rest/pkg/models"
)

func ValidateApiLoginPost(c *fiber.Ctx) error {
	var login models.Login
	if err := c.BodyParser(&login); err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid Credentials",
		})
	}
	v := validate.Struct(login)
	if !v.Validate() {
		return c.Status(401).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid Credentials",
		})
	}
	user, err := login.CheckLogin() //nolint:wsl

	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid Credentials",
		})
	}
	c.Locals("user", user)
	return c.Next()
}
