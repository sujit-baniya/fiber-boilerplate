package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	. "github.com/itsursujit/fiber-boilerplate/app"
	"github.com/itsursujit/fiber-boilerplate/auth"
	"github.com/itsursujit/fiber-boilerplate/config"
	"github.com/itsursujit/fiber-boilerplate/libraries"
	"github.com/itsursujit/fiber-boilerplate/models"
	"time"
)

func RegisterGet(c *fiber.Ctx) error {
	return c.Render("auth/register", fiber.Map{"Title": "Register"}, "layouts/auth")
}

func RegisterPost(c *fiber.Ctx) error {
	register := c.Locals("register").(models.RegisterForm)
	user, err := register.Signup()
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": true, "message": "Error on register request", "data": err.Error()}) //nolint:errcheck
	}
	store := Session.Get(c)       // get/create new session
	store.Set("user_id", user.ID) // save to storage
	_ = store.Save()

	go SendConfirmationEmail(user.Email, c.BaseURL())
	_ = c.JSON(user)
	return c.Redirect("/")
}

func VerifyRegisteredEmail(c *fiber.Ctx) error {
	return c.Redirect("/")
}

func ResendConfirmEmail(c *fiber.Ctx) error {
	user, _ := auth.User(c)
	go SendConfirmationEmail(user.Email, c.BaseURL())
	return c.Redirect("/")
}

func SendPasswordResetEmail(email string, baseURL string) {
	resetEmail := fmt.Sprintf("%s-reset-%d", email, time.Now().Unix())
	resetLink := GeneratePasswordResetURL(resetEmail, baseURL)
	htmlBody := config.PrepareHtml("emails/password-reset", fiber.Map{
		"reset_link": resetLink,
	})
	config.Send(email, "You asked to reset? Please click here!", htmlBody, "", "")
}

func RequestPasswordReset(c *fiber.Ctx) error {
	return c.Render("auth/request-password-reset", fiber.Map{"Title": "Reset Password"}, "layouts/auth")
}

func RequestPasswordResetPost(c *fiber.Ctx) error {
	email := c.FormValue("email")
	_, err := models.GetUserByEmail(email)
	if err != nil {
		return c.Redirect("/request-password-reset")
	}
	go SendPasswordResetEmail(email, c.BaseURL())
	return nil
}

func PasswordReset(c *fiber.Ctx) error {
	token := c.Query("t")
	return c.Render("auth/password-reset", fiber.Map{
		"Title": "Password Reset",
		"Token": token,
	}, "layouts/auth")
}

func PasswordResetPost(c *fiber.Ctx) error {
	register := c.Locals("register").(models.RegisterForm)
	email := c.Locals("email").(string)
	user, err := models.GetUserByEmail(email)
	if err != nil {
		return c.Status(401).SendString("Invalid Password Reset Token")
	}
	register.ID = user.ID
	_, err = register.ResetPassword()
	if err != nil {
		return c.Status(400).SendString("Oops!! Can't update password at the moment")
	}
	return c.Redirect("/login")
}

func SendConfirmationEmail(email string, baseURL string) {
	confirmLink := GenerateConfirmURL(email, baseURL)
	htmlBody := config.PrepareHtml("emails/confirm", fiber.Map{
		"confirm_link": confirmLink,
	})
	config.Send(email, "Is it you? Please confirm!", htmlBody, "", "")
}

func GenerateConfirmURL(email string, baseURL string) string {
	token := libraries.Encrypt(email, config.AppConfig.App_Key)
	uri := fmt.Sprintf("%s/do/verify-email?t=%s", baseURL, token)
	return uri
}

func GeneratePasswordResetURL(email string, baseURL string) string {
	token := libraries.Encrypt(email, config.AppConfig.App_Key)
	uri := fmt.Sprintf("%s/reset-password?t=%s", baseURL, token)
	return uri
}
