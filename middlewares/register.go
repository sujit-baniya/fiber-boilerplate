package middlewares

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"
	. "github.com/itsursujit/fiber-boilerplate/app"
	"github.com/itsursujit/fiber-boilerplate/auth"
	"github.com/itsursujit/fiber-boilerplate/config"
	"github.com/itsursujit/fiber-boilerplate/libraries"
	"github.com/itsursujit/fiber-boilerplate/models"
)

func ValidateRegisterPost(c *fiber.Ctx) error {
	var register models.RegisterForm
	if err := c.BodyParser(&register); err != nil {
		fmt.Println(err)
		return Flash.WithError(c, fiber.Map{
			"message": err.Error(),
		}).Redirect("/register")
	}

	v := validate.Struct(register)
	if !v.Validate() {
		fmt.Println(v.Errors)
		return Flash.WithError(c, fiber.Map{
			"message": v.Errors.One(),
		}).Redirect("/register")
	}
	c.Locals("register", register)
	return c.Next()
}

func ValidateConfirmToken(c *fiber.Ctx) error {
	t := libraries.Decrypt(c.Query("t"), config.AppConfig.App_Key)
	user, err := models.GetUserByEmail(t)
	if err != nil {
		return Flash.WithError(c, fiber.Map{
			"message": err.Error(),
		}).Redirect("/login")
	}

	if user.EmailVerified {
		return Flash.WithError(c, fiber.Map{
			"message": "Email was already validated",
		}).Redirect("/login")
	}
	user.EmailVerified = true
	DB.Save(&user)
	auth.Login(c, user.ID, config.AuthConfig.App_Jwt_Secret) //nolint:wsl
	c.Locals("user", user)
	return c.Next()
}
