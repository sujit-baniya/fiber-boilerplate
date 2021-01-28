package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sujit-baniya/verify-rest/app"
	"github.com/sujit-baniya/verify-rest/pkg/auth"
	"github.com/sujit-baniya/verify-rest/pkg/models"
	"github.com/sujit-baniya/verify-rest/pkg/services"
)

func RegisterGet(c *fiber.Ctx) error {
	data := getFlashData(c)
	data["title"] = "Register | "
	if err := c.Render("auth/register", data, "layouts/landing"); err != nil { //nolint:wsl
		panic(err.Error())
	}
	return nil
}

func RegisterPost(c *fiber.Ctx) error {
	register := c.Locals("register").(models.RegisterForm)
	user, err := register.Signup()
	if err != nil {
		return app.Http.Flash.WithError(c, fiber.Map{
			"error":   true,
			"message": "Error on register request: " + err.Error(),
		}).Redirect("/register")
	}
	store := app.Http.Session.Get(c) // get/create new session
	store.Set("user_id", user.ID)    // save to storage
	_ = store.Save()

	go services.SendConfirmationEmail(user.Email, app.Http.Server.Url)
	return c.Redirect("/")
}

func VerifyRegisteredEmail(c *fiber.Ctx) error {
	return c.Redirect("/app")
}

func ResendConfirmEmail(c *fiber.Ctx) error {
	user, _ := auth.User(c)
	go services.SendConfirmationEmail(user.Email, app.Http.Server.Url)
	return c.Redirect("/")
}

func RequestPasswordResetPost(c *fiber.Ctx) error {
	email := c.FormValue("email")
	_, err := models.GetUserByEmail(email)
	if err != nil {
		ctx := app.Http.Flash.WithError(c, fiber.Map{
			"message": "User with requested email doesn't exist",
		})
		fmt.Println(app.Http.Flash.Data)
		return ctx.Redirect("/request-password-reset")
	}
	go services.SendPasswordResetEmail(email, app.Http.Server.Url)
	return app.Http.Flash.WithSuccess(c, fiber.Map{
		"message": "We've sent an email for password to your registered email address",
	}).Redirect("/request-password-reset")
}

func PasswordReset(c *fiber.Ctx) error {
	token := c.Query("t")
	if err := c.Render("auth/password-reset", fiber.Map{
		"Title": "Password Reset",
		"Token": token,
	}, "layouts/landing"); err != nil { //nolint:wsl
		panic(err.Error())
	}
	return nil
}

func PasswordResetPost(c *fiber.Ctx) error {
	register := c.Locals("register").(models.RegisterForm)
	email := c.Locals("email").(string)
	user, err := models.GetUserByEmail(email)
	if err != nil {
		c.Status(401)
		return c.Send([]byte("Invalid Password Reset Token"))
	}
	register.ID = user.ID
	_, err = register.ResetPassword()
	if err != nil {
		return c.Status(400).Send([]byte("Oops!! Can't update password at the moment"))
	}
	return c.Redirect("/login")
}

func RequestPasswordReset(c *fiber.Ctx) error {
	data := getFlashData(c)
	data["title"] = "Password Reset | "
	if err := c.Render("auth/request-password-reset", data, "layouts/landing"); err != nil { //nolint:wsl
		panic(err.Error())
	}
	return nil
}
