package controllers

import (
	"github.com/gofiber/fiber/v2" //nolint:goimports
	"github.com/gookit/validate"
	"github.com/sujit-baniya/fiber-boilerplate/auth"
	"github.com/sujit-baniya/fiber-boilerplate/config"
	"github.com/sujit-baniya/fiber-boilerplate/models"
)

func OAuthToken(c *fiber.Ctx) error { //nolint:wsl
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
			"message": v.Errors.All(),
		})
	}
	user, err := login.CheckLogin() //nolint:wsl
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}
	token, err := auth.Login(c, user.ID, config.AuthConfig.Api_Jwt_Secret) //nolint:wsl
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"token":      token.Hash,
		"expires_in": token.Expire,
	})
}

func RefreshOauthToken(c *fiber.Ctx) {

}
