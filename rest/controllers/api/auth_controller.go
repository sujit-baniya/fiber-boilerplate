package controllers

import (
	"github.com/gofiber/fiber/v2" //nolint:goimports
	"github.com/gookit/validate"
	"github.com/sujit-baniya/fiber-boilerplate/app"
	"github.com/sujit-baniya/fiber-boilerplate/pkg/auth"
	"github.com/sujit-baniya/fiber-boilerplate/pkg/models"
	"github.com/sujit-baniya/fiber-boilerplate/pkg/services"
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
	token, err := auth.Login(c, user.ID, app.Http.Token.ApiJwtSecret) //nolint:wsl
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

func ApiLoginPost(c *fiber.Ctx) error { //nolint:wsl
	user := c.Locals("user").(*models.User)
	auth.Login(c, user.ID, app.Http.Token.AppJwtSecret) //nolint:wsl
	return c.JSON(user)
}

func ApiRegisterPost(c *fiber.Ctx) error {
	register := c.Locals("register").(models.RegisterForm)
	user, err := register.Signup()
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": true, "message": "Error on register request", "data": err.Error()}) //nolint:errcheck

	}
	store := app.Http.Session.Get(c) // get/create new session
	store.Set("user_id", user.ID)    // save to storage
	_ = store.Save()

	go services.SendConfirmationEmail(user.Email, c.BaseURL())
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Registered successfully! Please confirm your email",
	})
}

func RefreshOauthToken(c *fiber.Ctx) error {
	return nil
}
